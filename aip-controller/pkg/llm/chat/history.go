package chat

import (
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
)

const ChatHistoryContextKey chain.ContextKey[Message] = "chat_history"
const ChatReplyContextKey chain.ContextKey[Message] = "chat_reply"

type ChatHistory interface {
	AsSlice() []MessageEntry
	AsMessage() Message
}
