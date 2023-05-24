package memory

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
)

type ChannelChatMemory struct {
	Messenger msn.API
	Channel   msn.ChannelID
	Endpoint  msn.EndpointID

	Database forddb.Database

	ContextKey chain.ContextKey[chat.Message]

	joined bool
}

func (i *ChannelChatMemory) Create(ctx context.Context) error {
	_, err := i.Messenger.JoinChannel(ctx, &msn.JoinChannelRequest{
		ChannelID:  i.Channel,
		EndpointID: i.Endpoint,
	})

	if err != nil {
		return err
	}

	return nil
}

func (i *ChannelChatMemory) Load(ctx chain.ChainContext) (chat.Message, error) {
	if !i.joined {
		err := i.Create(ctx.Context())

		if err != nil {
			return chat.Message{}, err
		}

		i.joined = true
	}

	messages, err := i.Database.List(
		ctx.Context(),
		msn.MessageID{}.BasicResourceType().GetResourceID(),
		forddb.WithSortField("metadata.created_at", forddb.Desc),
		forddb.WithListQueryOptions(
			forddb.WithFilterParameter("channel", i.Channel.String()),
			forddb.WithFilterExpression(`resource.channel_id == args.channel`),
		),
	)

	if err != nil {
		return chat.Message{}, err
	}

	msg := chat.Message{}
	msg.Entries = make([]chat.MessageEntry, len(messages))

	for i, message := range messages {
		msg.Entries[i] = *(*chat.MessageEntry)(message.(*msn.Message))
	}

	return msg, nil
}

func (i *ChannelChatMemory) Append(ctx chain.ChainContext, msg chat.Message) error {
	if !i.joined {
		err := i.Create(ctx.Context())

		if err != nil {
			return err
		}

		i.joined = true
	}

	for _, entry := range msg.Entries {
		req := &msn.PostMessageRequest{
			From:    entry.From,
			Channel: entry.ChannelID,
			Text:    entry.Text,
		}

		if req.From.IsEmpty() {
			req.From = i.Endpoint
		}

		if req.Channel.IsEmpty() {
			req.Channel = i.Channel
		}

		_, err := i.Messenger.PostMessage(ctx.Context(), req)

		if err != nil {
			return err
		}
	}

	return nil
}
