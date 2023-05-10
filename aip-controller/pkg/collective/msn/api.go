package msn

import (
	"context"
)

type PostMessageRequest struct {
	From    EndpointID `json:"from"`
	Channel ChannelID  `json:"channel"`

	Text string `json:"text"`
}

type PostMessageResponse struct {
	Message Message `json:"message"`
}

type JoinChannelRequest struct {
	EndpointID EndpointID `json:"endpoint_id"`
	ChannelID  ChannelID  `json:"channel_id"`
}

type JoinChannelResponse struct {
	Channel Channel `json:"channel"`
}

type LeaveChannelRequest struct {
	EndpointID EndpointID `json:"endpoint_id"`
	ChannelID  ChannelID  `json:"channel_id"`
}

type LeaveChannelResponse struct {
	Channel Channel `json:"channel"`
}

type SubscribeToEventsRequest struct {
	EndpointID EndpointID `json:"endpoint_id"`
}

type MessageEvent struct {
	Message Message `json:"message"`
}

type Event struct {
	MessageEvent *MessageEvent `json:"message_event"`
}

type API interface {
	// JoinChannel joins a channel.
	JoinChannel(ctx context.Context, req *JoinChannelRequest) (*JoinChannelResponse, error)

	// LeaveChannel leaves a channel.
	LeaveChannel(ctx context.Context, req *LeaveChannelRequest) (*LeaveChannelResponse, error)

	// PostMessage sends a message to a channel.
	PostMessage(ctx context.Context, req *PostMessageRequest) (*PostMessageResponse, error)

	// SubscribeToEvents subscribes to a channel.
	SubscribeToEvents(ctx context.Context, req *SubscribeToEventsRequest) (<-chan Event, error)
}
