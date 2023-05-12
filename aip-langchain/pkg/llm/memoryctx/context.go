package memoryctx

import (
	"context"

	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
)

const memoryKey = "llm-chat-memory"

func WithMemory(ctx context.Context, memory chat.Memory) context.Context {
	return context.WithValue(ctx, memoryKey, memory)
}

func GetMemory(ctx context.Context) chat.Memory {
	return ctx.Value(memoryKey).(chat.Memory)
}
