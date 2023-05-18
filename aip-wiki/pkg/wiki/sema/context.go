package sema

import (
	context "context"

	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/chunkers"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore/faiss"
)

type SemanticContext struct {
	Tokenizer tokenizers.BasicTokenizer
	Chunker   chunkers.Chunker

	Model    chat.LanguageModel
	Embedder llm.Embedder

	SummarizeChain chain.Chain
	CompressChain  chain.Chain
	RefineChain    chain.Chain

	Index vectorstore.Index
}

const ContextKey chain.ContextKey[string] = "Context"
const DocumentsKey chain.ContextKey[[]*SemanticNode] = "Documents"
const CompressionRatioKey chain.ContextKey[int] = "CompressionRatio"

var SummarizePrompt = chat.ComposeTemplate(
	chat.EntryTemplate(msn.RoleSystem,
		chain.NewTemplatePrompt(
			`
You are an AI assistant specialized in comparing, summarizing, merging and splitting documents.
Perform the action request by the user for the documents below.

Summarize the documents below into a section-based tree essay as single cohesive document.`,
		),
	),

	chat.EntryTemplate(msn.RoleUser,
		chain.NewTemplatePrompt(`
{{ range $index, $doc := .Documents }}
<<< Document {{ $index }} >>>
{{ $doc.Status.Value }}
{{ end }}
`,
			chain.WithRequiredInput(DocumentsKey)),
	),
)

var CompressPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(msn.RoleSystem,
		chain.NewTemplatePrompt(
			`
You are an AI assistant specialized in comparing, summarizing, merging and splitting documents.
Perform the action request by the user for the documents below.

Summarize the following documents with a {{ .CompressionRatio }}:1 compression ratio, to the best of your ability.
`,
			chain.WithRequiredInput(CompressionRatioKey),
		),
	),

	chat.EntryTemplate(msn.RoleUser,
		chain.NewTemplatePrompt(`
{{ range $index, $doc := .Documents }}
<<< Document {{ $index }} >>>
{{ $doc.Status.Value }}
{{ end }}
`,
			chain.WithRequiredInput(DocumentsKey)),
	),
)

var RefinePrompt = chat.ComposeTemplate(
	chat.EntryTemplate(msn.RoleSystem,
		chain.NewTemplatePrompt(
			`
You are an AI assistant specialized in comparing, summarizing, merging and splitting documents.
Perform the action request by the user for the documents below.

Summarize the documents below into a section-based tree essay as single cohesive document based on the context below.

<<< Context >>>
{{ range $index, $doc := .Context }}
{{ $doc.Status.Value }}
{{ end }}
`,
			chain.WithRequiredInput(ContextKey),
		),
	),

	chat.EntryTemplate(msn.RoleUser,
		chain.NewTemplatePrompt(`
{{ range $index, $doc := .Documents }}
<<< Document {{ $index }} >>>
{{ $doc.Status.Value }}
{{ end }}
`,
			chain.WithRequiredInput(DocumentsKey)),
	),
)

var SplitBySectionPrompt = chat.ComposeTemplate(
	chat.EntryTemplate(msn.RoleSystem,
		chain.NewTemplatePrompt(
			`
You are an AI assistant specialized in comparing, summarizing, merging and splitting documents.
Perform the action request by the user for the documents below.

Summarize the documents below into a section-based tree as single cohesive document.
`,
			chain.WithRequiredInput(CompressionRatioKey),
		),
	),

	chat.EntryTemplate(msn.RoleUser,
		chain.NewTemplatePrompt(`
{{ range $index, $doc := .Documents }}
<<< Document {{ $index }} >>>
{{ $doc.Status.Value }}
{{ end }}
`,
			chain.WithRequiredInput(DocumentsKey)),
	),
)

func InitializeSemanticContext(sc *SemanticContext) {
	var err error

	sc.Index, err = faiss.NewIndex(1536)

	if err != nil {
		panic(err)
	}

	sc.SummarizeChain = chain.New(
		chain.Sequential(
			chat.Predict(sc.Model, SummarizePrompt),
		),
	)

	sc.CompressChain = chain.New(
		chain.Sequential(
			chat.Predict(sc.Model, CompressPrompt),
		),
	)

	sc.RefineChain = chain.New(
		chain.Sequential(
			chat.Predict(sc.Model, RefinePrompt),
		),
	)
}

func (sc *SemanticContext) Content(content string) SemanticContent {
	count, err := sc.Tokenizer.Count(content)

	if err != nil {
		panic(err)
	}

	return SemanticContent{
		Content:    content,
		TokenCount: count,
	}
}

func (sc *SemanticContext) Unit(content ...SemanticContent) SemanticUnit {
	return SemanticUnit{Contents: content}
}

func (sc *SemanticContext) Append(ctx context.Context, node *SemanticNode) error {
	doc := vectorstore.Document{}
	doc.Content = node.Status.Value.String()
	doc.Type = "SemanticNode"
	doc.ID = node.ID.String()

	_, err := sc.Index.IndexDocument(ctx, &doc, vectorstore.WithIndexEmbedder(sc.Embedder))

	if err != nil {
		return err
	}

	return nil
}

func (sc *SemanticContext) Recall(ctx context.Context, query string) (*SemanticNode, error) {
	res, err := sc.Index.Search(
		ctx,
		query,
		vectorstore.WithSearchEmbedder(sc.Embedder),
		vectorstore.WithReturnHitContents(),
	)

	if err != nil {
		return nil, err
	}

	result := &SemanticNode{}

	for _, v := range res.Hits {
		unit := sc.Unit(sc.Content(v.Content))
		hitId := forddb.NewStringID[SemanticNodeID](v.ID)

		result.Spec.ParentNodeIds = append(result.Spec.ParentNodeIds, hitId)
		result.Status.Value = result.Status.Value.Add(unit)
	}

	result.Status.Value, err = result.Status.Value.Reshape(ctx, sc.Chunker, 300)

	if err != nil {
		return nil, err
	}

	result, err = sc.Merge(ctx, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sc *SemanticContext) Refine(
	ctx context.Context,
	root *SemanticNode,
) (result *SemanticNode, _ error) {
	contextNode, err := sc.Recall(ctx, root.Status.Value.String())

	if err != nil {
		return result, err
	}

	if err := sc.Append(ctx, contextNode); err != nil {
		return nil, err
	}

	contextNode, err = sc.Mip(ctx, contextNode)

	if err != nil {
		return result, err
	}

	result = &SemanticNode{}
	result.Spec.SemanticRootID = root.ID
	result.Spec.SemanticRootHeight = root.Spec.SemanticRootHeight + 1
	result.Spec.SemanticRootRow = 0
	result.Spec.ParentNodeIds = []SemanticNodeID{root.ID}

	cctx := chain.NewChainContext(ctx)
	cctx.SetInput(ContextKey, []*SemanticNode{contextNode})
	cctx.SetInput(DocumentsKey, []*SemanticNode{root})

	if err := sc.RefineChain.Run(cctx); err != nil {
		return result, err
	}

	reply := chain.Output(cctx, chat.ChatReplyContextKey)
	output := reply.AsText()

	result.Status.Value = sc.Unit(sc.Content(output))

	return
}

func (sc *SemanticContext) Merge(
	ctx context.Context,
	root *SemanticNode,
	nodes ...*SemanticNode,
) (result *SemanticNode, _ error) {
	allNodes := append([]*SemanticNode{root}, nodes...)

	result = &SemanticNode{}
	result.Spec.SemanticRootID = root.ID
	result.Spec.SemanticRootHeight = root.Spec.SemanticRootHeight + 1
	result.Spec.SemanticRootRow = 0

	result.Spec.ParentNodeIds = lo.Map(allNodes, func(item *SemanticNode, index int) SemanticNodeID {
		return item.ID
	})

	cctx := chain.NewChainContext(ctx)
	cctx.SetInput(DocumentsKey, allNodes)

	if err := sc.SummarizeChain.Run(cctx); err != nil {
		return result, err
	}

	reply := chain.Output(cctx, chat.ChatReplyContextKey)
	output := reply.AsText()

	result.Status.Value = sc.Unit(sc.Content(output))

	return
}

func (sc *SemanticContext) Mip(
	ctx context.Context,
	root *SemanticNode,
) (result *SemanticNode, _ error) {
	result = &SemanticNode{}
	result.Spec.SemanticRootID = root.ID
	result.Spec.SemanticRootHeight = root.Spec.SemanticRootHeight - 1
	result.Spec.SemanticRootRow = 0

	result.Spec.ParentNodeIds = []SemanticNodeID{root.ID}

	cctx := chain.NewChainContext(ctx)
	cctx.SetInput(DocumentsKey, []*SemanticNode{root})
	cctx.SetInput(CompressionRatioKey, 2)

	if err := sc.CompressChain.Run(cctx); err != nil {
		return result, err
	}

	reply := chain.Output(cctx, chat.ChatReplyContextKey)
	output := reply.AsText()

	result.Status.Value = sc.Unit(sc.Content(output))

	return
}
