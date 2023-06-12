package mdcodegen

import (
	"context"
	"flag"
	"image"
	"image/color"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget/material"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/cobra"
	"github.com/zyedidia/generic/stack"
	"go.uber.org/zap"

	"gioui.org/font/opentype"
	"github.com/steverusso/gio-fonts/inconsolata/inconsolatabold"
	"github.com/steverusso/gio-fonts/inconsolata/inconsolataregular"
	"github.com/steverusso/gio-fonts/nunito/nunitobold"
	"github.com/steverusso/gio-fonts/nunito/nunitobolditalic"
	"github.com/steverusso/gio-fonts/nunito/nunitoitalic"
	"github.com/steverusso/gio-fonts/nunito/nunitoregular"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/memory"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
	"github.com/greenboxal/aip/aip-wiki/exp/mdedit"
)

var DocumentKey chain.ContextKey[string] = "Document"

var CodeGeneratorPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(
		msn.RoleSystem,
		chain.NewTemplatePrompt(`
You are an AI agent specialized in generating code in Go.
You are given a Document containing code.

### Document
`+"```go"+`
{{.Document}}
`+"```"+`
	`, chain.WithRequiredInput(DocumentKey)),
	),

	chat.HistoryFromContext(memory.ContextualMemoryKey),

	chat.EntryTemplate(
		msn.RoleUser,
		chain.NewTemplatePrompt(
			`Complete the TODOs in the document above.`,
		),
	),

	chat.EntryTemplate(msn.RoleAI, chain.Static("")),
)

type MarkdownCodeGenCommand struct {
	logger *zap.SugaredLogger

	oai      *openai.Client
	embedder *openai.Embedder

	contentChain chain.Chain
	model        *openai.ChatLanguageModel

	fsys        diskFS
	editorReady chan struct{}
	editorEvent chan interface{}
}

func NewMarkdownCodeGenCommand(
	logger *zap.SugaredLogger,
	oai *openai.Client,
) *MarkdownCodeGenCommand {
	mcgen := &MarkdownCodeGenCommand{
		logger: logger,
		oai:    oai,

		editorReady: make(chan struct{}),
		editorEvent: make(chan interface{}, 10),
	}

	mcgen.embedder = &openai.Embedder{
		Client: oai,
		Model:  openai.AdaEmbeddingV2,
	}

	mcgen.model = &openai.ChatLanguageModel{
		Client:      oai,
		Model:       "gpt-3.5-turbo",
		Temperature: 0.0,
	}

	mcgen.contentChain = chain.New(
		chain.WithName("MarkdownCodeGenerator"),

		chain.Sequential(
			chat.Predict(
				mcgen.model,
				CodeGeneratorPrompt,
			),
		),
	)

	return mcgen
}

func (mcgen *MarkdownCodeGenCommand) Run(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	mcgen.fsys.document = fileBytes

	documentRoot := ParseMarkdown(fileBytes)

	if err := mcgen.ProcessDocument(filePath, documentRoot); err != nil {
		return err
	}

	mcgen.logger.Info("Done")

	return nil
}

type MarkdownCodeGenContext struct {
	stack   *stack.Stack[MarkdownCodeGenState]
	builder *CompilationUnitBuilder

	State MarkdownCodeGenState
}

func NewMarkdownCodeGenContext() *MarkdownCodeGenContext {
	return &MarkdownCodeGenContext{
		stack:   stack.New[MarkdownCodeGenState](),
		builder: NewCompilationUnitBuilder(),
	}
}

func (mcgen *MarkdownCodeGenContext) PushState() {
	mcgen.stack.Push(mcgen.State)
}

func (mcgen *MarkdownCodeGenContext) PopState() {
	mcgen.State = mcgen.stack.Pop()
}

type MarkdownCodeGenState struct {
	CurrentPackage  string
	CurrentFile     string
	CurrentLanguage string
	CurrentCode     string

	CurrentModule *ModuleBuilder
}

func (mcgen *MarkdownCodeGenCommand) ProcessNode(root ast.Node) (isFinal bool, err error) {
	ctx := NewMarkdownCodeGenContext()

	isFinal = true

	ast.WalkFunc(root, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch node := node.(type) {
			case *ast.Heading:
				name := string(node.Literal)

				if len(node.Children) > 0 {
					name = string(node.Children[0].AsLeaf().Literal)
				}

				if strings.HasPrefix(name, "@file: ") {
					ctx.State.CurrentFile = strings.TrimPrefix(name, "@file: ")
				} else if strings.HasPrefix(name, "@package: ") {
					ctx.State.CurrentPackage = strings.TrimPrefix(name, "@package: ")
					ctx.State.CurrentModule = ctx.builder.GetOrCreateModule(ctx.State.CurrentPackage)
				}

			case *ast.CodeBlock:
				ctx.State.CurrentLanguage = string(node.Info)
				ctx.State.CurrentCode = string(node.Literal)

				ctx.State.CurrentModule.AddFile(ctx.State.CurrentFile, ctx.State.CurrentCode)

				/*updated, e := mcgen.ProcessCodeBlock(
					ctx.State.CurrentPackage,
					ctx.State.CurrentFile,
					ctx.State.CurrentLanguage,
					ctx.State.CurrentCode,
				)

				if e != nil {
					err = e
					return ast.Terminate
				}*

				node.Literal = []byte(updated)*/
			}
		} else {
		}

		return ast.GoToNext
	})

	if err != nil {
		return
	}

	cu, err := ctx.builder.Build()

	if err != nil {
		return
	}

	cu = cu

	return isFinal, nil
}

func (mcgen *MarkdownCodeGenCommand) ProcessDocument(path string, root ast.Node) error {
	currentRound := 0
	outputPath := strings.Replace(path, ".md", ".out.md", 1)

	for {
		currentRound++

		mcgen.logger.Infow("Beginning processing round", "round", currentRound)

		isFinal, err := mcgen.ProcessNode(root)

		if err != nil {
			return err
		}

		result := FormatMarkdown(root)

		mcgen.fsys.document = result

		if err := os.WriteFile(outputPath, result, 0644); err != nil {
			return err
		}

		if isFinal {
			break
		}
	}

	mcgen.logger.Infow("Finished processing document", "rounds", currentRound)

	return nil
}

func (mcgen *MarkdownCodeGenCommand) ProcessCodeBlock(
	pkg string,
	file string,
	language string,
	code string,
) (string, error) {
	mcgen.logger.Infow("Processing code block", "pkg", pkg, "file", file, "language", language)

	ctx := context.Background()
	cctx := chain.NewChainContext(ctx)

	cctx.SetInput(DocumentKey, code)

	if err := mcgen.contentChain.Run(cctx); err != nil {
		return "", err
	}

	result := chain.Output(cctx, chat.ChatReplyContextKey)
	reply := result.Entries[0].Text
	codeOutput := ""
	parsedReply := ParseMarkdown([]byte(reply))

	ast.WalkFunc(parsedReply, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch node := node.(type) {
			case *ast.CodeBlock:
				codeOutput += string(node.Literal)
			}
		}

		return ast.GoToNext
	})

	return codeOutput, nil
}

func (mcgen *MarkdownCodeGenCommand) runDebugger() error {
	var ops op.Ops

	win := app.NewWindow(
		app.Size(800, 600),
		app.Title("MdEdit"),
	)

	win.Perform(system.ActionCenter)

	th := material.NewTheme([]font.FontFace{
		// Proportionals.
		mustFont(font.Font{}, nunitoregular.TTF),
		mustFont(font.Font{Weight: font.Bold}, nunitobold.TTF),
		mustFont(font.Font{Weight: font.Bold, Style: font.Italic}, nunitobolditalic.TTF),
		mustFont(font.Font{Style: font.Italic}, nunitoitalic.TTF),
		// Monos.
		mustFont(font.Font{Variant: "Mono"}, inconsolataregular.TTF),
		mustFont(font.Font{Variant: "Mono", Weight: font.Bold}, inconsolatabold.TTF),
	})
	th.TextSize = 18
	th.Palette = material.Palette{
		Bg:         color.NRGBA{17, 21, 24, 255},
		Fg:         color.NRGBA{235, 235, 235, 255},
		ContrastFg: color.NRGBA{10, 180, 230, 255},
		ContrastBg: color.NRGBA{220, 220, 220, 255},
	}

	s := mdedit.NewSession(&mcgen.fsys, win)
	for _, fpath := range flag.Args() {
		s.OpenFile(fpath)
	}
	s.FocusActiveTab()

	close(mcgen.editorReady)

	for {
		select {
		case e := <-win.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				start := time.Now()
				gtx := layout.NewContext(&ops, e)
				// Process any key events since the previous frame.
				for _, ke := range gtx.Events(win) {
					if ke, ok := ke.(key.Event); ok {
						s.HandleKeyEvent(ke)
					}
				}
				// Gather key input on the entire window area.
				areaStack := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
				key.InputOp{Tag: win, Keys: topLevelKeySet}.Add(gtx.Ops)
				s.Layout(gtx, th)
				areaStack.Pop()

				e.Frame(gtx.Ops)
				if *printFrameTimes {
					log.Println(time.Since(start))
				}

				for done := false; !done; {
					select {
					case <-mcgen.editorEvent:
						s.OpenFile("/document")
					default:
						done = true
					}
				}

			case system.DestroyEvent:
				return e.Err
			}
		}
	}
}

func FormatMarkdown(node ast.Node) []byte {
	return markdown.Render(node, md.NewRenderer())
}

func ParseMarkdown(md []byte) ast.Node {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	return p.Parse(md)
}

const topLevelKeySet = "Ctrl-[O,W," + key.NameTab + "]" +
	"|Ctrl-Shift-[" + key.NamePageUp + "," + key.NamePageDown + "," + key.NameTab + "]" +
	"|Alt-[1,2,3,4,5,6,7,8,9]"

var printFrameTimes = flag.Bool("print-frame-times", false, "Print how long each frame takes.")

type diskFS struct {
	document []byte
}

func newDiskFS() (*diskFS, error) {
	return &diskFS{}, nil
}

func (d *diskFS) HomeDir() string {
	return "/"
}

func (d *diskFS) WorkingDir() string {
	return "/"
}

func (*diskFS) ReadDir(fpath string) ([]fs.FileInfo, error) {
	return nil, nil
}

func (fs *diskFS) ReadFile(fpath string) ([]byte, error) {
	if fs.document == nil {
		return nil, nil
	}

	return fs.document, nil
}

func (*diskFS) WriteFile(fpath string, data []byte) error {
	return nil
}

func mustFont(fnt font.Font, data []byte) font.FontFace {
	face, err := opentype.Parse(data)
	if err != nil {
		panic("failed to parse font: " + err.Error())
	}
	return font.FontFace{Font: fnt, Face: face}
}
