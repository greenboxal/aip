package msn

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type MessageID struct {
	forddb.StringResourceID[*Message] `ipld:",inline"`
}

type Message struct {
	forddb.ResourceBase[MessageID, *Message] `json:"metadata"`

	ThreadID  string `json:"thread_id"`
	ReplyToID string `json:"reply_to_id"`

	Channel ChannelID  `json:"channel"`
	From    EndpointID `json:"from"`
	Name    string     `json:"username"`
	Role    Role       `json:"role"`
	Text    string     `json:"text"`
}
