package memory

import (
	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
	chat2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

const ContextualMemoryKey chain.ContextKey[chat2.Message] = "contextual_memory"

type ContextualMemory struct {
	HistoryKey chain.ContextKey[chat2.Memory]
	InputKey   chain.ContextKey[chat2.Message]
	ContextKey chain.ContextKey[chat2.Message]

	Index    vectorstore.VectorStore
	Embedder llm.Embedder
}

func (cm *ContextualMemory) Load(ctx chain.ChainContext) error {
	//chatMemory := chain.Input(ctx, cm.HistoryKey)
	currentInput := chain.Input(ctx, cm.InputKey)

	return cm.LoadFor(ctx, currentInput.AsText())
}

func (cm *ContextualMemory) LoadFor(ctx chain.ChainContext, query string) error {
	result, err := cm.Index.Search(
		ctx.Context(),
		query,
		vectorstore.WithReturnHitContents(),
		vectorstore.WithSearchEmbedder(cm.Embedder),
	)

	if err != nil {
		return err
	}

	if len(result.Hits) == 0 {
		chain.SetOutput(ctx, cm.ContextKey, chat2.Message{})
		return nil
	}

	resultAsText := result.Hits[0].Content

	msg := chat2.Compose(
		chat2.Entry(msn.RoleSystem, "Contextual memory related to the request: "+resultAsText),
	)

	chain.SetOutput(ctx, cm.ContextKey, msg)

	return nil
}
