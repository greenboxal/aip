package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/greenboxal/aip/aip-controller/pkg/collective"
)

type Transport struct {
	ps *pubsub.PubSub

	m          sync.Mutex
	topics     map[string]*pubsub.Topic
	incomingCh chan collective.Message
}

func NewTransport(ps *pubsub.PubSub) *Transport {
	return &Transport{
		ps: ps,

		topics:     map[string]*pubsub.Topic{},
		incomingCh: make(chan collective.Message, 16),
	}
}

func (t *Transport) getTopic(channel string, create bool) (*pubsub.Topic, error) {
	t.m.Lock()
	defer t.m.Unlock()

	topic := t.topics[channel]

	if topic != nil {
		return topic, nil
	}

	if !create {
		return nil, nil
	}

	topic, err := t.ps.Join(fmt.Sprintf("aip-ps-%s", channel))

	if err != nil {
		return nil, err
	}

	t.topics[channel] = topic

	return topic, nil
}

func (t *Transport) Subscribe(channel string) error {
	_, err := t.getTopic(channel, true)

	if err != nil {
		return err
	}

	return nil
}

func (t *Transport) Incoming() <-chan collective.Message {
	return t.incomingCh
}

func (t *Transport) RouteMessage(ctx context.Context, msg collective.Message) error {
	topic, err := t.getTopic(msg.Channel, false)

	if err != nil {
		return err
	}

	if topic == nil {
		return nil
	}

	data, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	return topic.Publish(ctx, data)
}

func (t *Transport) Close() error {
	return nil
}
