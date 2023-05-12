package msn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type Service struct {
	db     forddb.Database
	logger *zap.SugaredLogger
	router *Router
}

func NewService(
	logger *zap.SugaredLogger,
	db forddb.Database,
	router *Router,
) *Service {
	return &Service{
		logger: logger.Named("msn-service"),
		db:     db,
		router: router,
	}
}

func (r *Service) JoinChannel(ctx context.Context, req *JoinChannelRequest) (*JoinChannelResponse, error) {
	endpoint, err := forddb.Get[*Endpoint](ctx, r.db, req.EndpointID)

	if err != nil && !forddb.IsNotFound(err) {
		return nil, err
	}

	channel, err := forddb.Get[*Channel](ctx, r.db, req.ChannelID)

	if err != nil && !forddb.IsNotFound(err) {
		return nil, err
	}

	if endpoint == nil {
		endpoint = &Endpoint{}
		endpoint.ID = req.EndpointID
	}

	if channel == nil {
		channel = &Channel{}
		channel.ID = req.ChannelID
	}

	endpoint.Subscriptions = append(endpoint.Subscriptions, req.ChannelID)
	channel.Subscribers = append(channel.Subscribers, req.EndpointID)

	channel, err = forddb.Put[*Channel](ctx, r.db, channel)

	if err != nil {
		return nil, err
	}

	endpoint, err = forddb.Put[*Endpoint](ctx, r.db, endpoint)

	if err != nil {
		return nil, err
	}

	return &JoinChannelResponse{
		Channel: *channel,
	}, nil
}

func (r *Service) LeaveChannel(ctx context.Context, req *LeaveChannelRequest) (*LeaveChannelResponse, error) {
	return nil, errors.New("not implemented yet")
}

func (r *Service) PostMessage(ctx context.Context, req *PostMessageRequest) (*PostMessageResponse, error) {
	var err error

	_, err = forddb.Get[*Endpoint](ctx, r.db, req.From)

	if err != nil {
		return nil, err
	}

	_, err = forddb.Get[*Channel](ctx, r.db, req.Channel)

	if err != nil {
		return nil, err
	}

	msg := &Message{}
	msg.ID = forddb.NewStringID[MessageID](fmt.Sprintf("%d", time.Now().UnixNano()))
	msg.Channel = req.Channel
	msg.From = req.From
	msg.Text = req.Text

	msg, err = forddb.Put[*Message](ctx, r.db, msg)

	if err != nil {
		return nil, err
	}

	r.router.Dispatch(Event{
		MessageEvent: &MessageEvent{
			Message: *msg,
		},
	})

	return &PostMessageResponse{
		Message: *msg,
	}, nil
}
