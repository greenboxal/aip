package chat

import "github.com/greenboxal/aip/pkg/llm"

const ChatHistoryContextKey llm.ContextKey[Message] = "chat_history"
const ChatReplyContextKey llm.ContextKey[Message] = "chat_reply"

type ChatHistory interface {
	AsSlice() []MessageEntry
	AsMessage() Message
}
