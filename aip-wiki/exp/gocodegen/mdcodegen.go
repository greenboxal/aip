package mdcodegen

import (
	"context"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/cobra"
	"github.com/zyedidia/generic/stack"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/memory"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
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

type GoCodeGenCommand struct {
	logger *zap.SugaredLogger

	oai      *openai.Client
	embedder *openai.Embedder

	contentChain chain.Chain
	model        *openai.ChatLanguageModel
}

func NewGoCodeGenCommand(
	logger *zap.SugaredLogger,
	oai *openai.Client,
) *GoCodeGenCommand {
	mcgen := &GoCodeGenCommand{
		logger: logger,
		oai:    oai,
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
		chain.WithName("GoCodeGenerator"),

		chain.Sequential(
			chat.Predict(
				mcgen.model,
				CodeGeneratorPrompt,
			),
		),
	)

	return mcgen
}

func (mcgen *GoCodeGenCommand) Run(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	documentRoot := ParseMarkdown(fileBytes)

	if err := mcgen.ProcessDocument(filePath, documentRoot); err != nil {
		return err
	}

	mcgen.logger.Info("Done")

	return nil
}

type GoCodeGenContext struct {
	stack   *stack.Stack[GoCodeGenState]
	builder *CompilationUnitBuilder

	State GoCodeGenState
}

func NewGoCodeGenContext() *GoCodeGenContext {
	return &GoCodeGenContext{
		stack: stack.New[GoCodeGenState](),
	}
}

func (mcgen *GoCodeGenContext) PushState() {
	mcgen.stack.Push(mcgen.State)
}

func (mcgen *GoCodeGenContext) PopState() {
	mcgen.State = mcgen.stack.Pop()
}

type GoCodeGenState struct {
	CurrentPackage  string
	CurrentFile     string
	CurrentLanguage string
	CurrentCode     string

	CurrentModule *ModuleBuilder
}

func (mcgen *GoCodeGenCommand) ProcessNode(root ast.Node) (isFinal bool, err error) {
	ctx := NewGoCodeGenContext()

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

				updated, e := mcgen.ProcessCodeBlock(
					ctx.State.CurrentPackage,
					ctx.State.CurrentFile,
					ctx.State.CurrentLanguage,
					ctx.State.CurrentCode,
				)

				if e != nil {
					err = e
					return ast.Terminate
				}

				node.Literal = []byte(updated)
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

func (mcgen *GoCodeGenCommand) ProcessDocument(path string, root ast.Node) error {
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

func (mcgen *GoCodeGenCommand) ProcessCodeBlock(
	pkg string,
	file string,
	language string,
	code string,
) (string, error) {
	mcgen.logger.Infow("Processing code block", "pkg", pkg, "file", file, "language", language)

	ctx := context.Background()
	cctx := chain.NewChainContext(ctx)

	// TODO: Parse golang AST from `code` and use that to generate the prompt
	// It should walk over every type and function and generate a prompt for it

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

func FormatMarkdown(node ast.Node) []byte {
	return markdown.Render(node, md.NewRenderer())
}

func ParseMarkdown(md []byte) ast.Node {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	return p.Parse(md)
}
