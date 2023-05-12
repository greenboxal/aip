package chat

import (
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
)

const ChatHistoryContextKey chain.ContextKey[Message] = "chat_history"
const ChatReplyContextKey chain.ContextKey[Message] = "chat_reply"
