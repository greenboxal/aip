package memory

import (
	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing"
	"github.com/greenboxal/aip/aip-controller/pkg/llm"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chat"
)

const ContextualMemoryKey chain.ContextKey[chat.Message] = "contextual_memory"

type ContextualMemory struct {
	HistoryKey chain.ContextKey[chat.Memory]
	InputKey   chain.ContextKey[chat.Message]
	ContextKey chain.ContextKey[chat.Message]

	Index    indexing.Provider
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
		indexing.WithReturnHitContents(),
		indexing.WithSearchEmbedder(cm.Embedder),
	)

	if err != nil {
		return err
	}

	if len(result.Hits) == 0 {
		chain.SetOutput(ctx, cm.ContextKey, chat.Message{})
		return nil
	}

	resultAsText := result.Hits[0].Content

	msg := chat.Compose(
		chat.Entry(msn.RoleSystem, "Contextual memory related to the request: "+resultAsText),
	)

	chain.SetOutput(ctx, cm.ContextKey, msg)

	return nil
}
