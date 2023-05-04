package ipfs

import (
	context "context"
	"sync"

	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/options"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type PubSubManager struct {
	m      sync.RWMutex
	ps     *pubsub.PubSub
	host   host.Host
	topics map[string]*pubsubTopic
}

func NewPubSubManager(h host.Host, ps *pubsub.PubSub) *PubSubManager {
	return &PubSubManager{
		host:   h,
		ps:     ps,
		topics: map[string]*pubsubTopic{},
	}
}

func (p *PubSubManager) Ls(ctx context.Context) ([]string, error) {
	return p.ps.GetTopics(), nil
}

func (p *PubSubManager) Peers(ctx context.Context, option ...options.PubSubPeersOption) ([]peer.ID, error) {
	opts, err := options.PubSubPeersOptions(option...)

	if err != nil {
		return nil, err
	}

	if opts.Topic != "" {
		return p.ps.ListPeers(opts.Topic), nil
	}

	return p.host.Peerstore().Peers(), nil
}

func (p *PubSubManager) Publish(ctx context.Context, s string, bytes []byte) error {
	topic, err := p.getOrJoinTopic(ctx, s)

	if err != nil {
		return err
	}

	return topic.topic.Publish(ctx, bytes)
}

func (p *PubSubManager) Subscribe(ctx context.Context, s string, option ...options.PubSubSubscribeOption) (iface.PubSubSubscription, error) {
	_, err := options.PubSubSubscribeOptions(option...)

	if err != nil {
		return nil, err
	}

	topic, err := p.getOrJoinTopic(ctx, s)

	if err != nil {
		return nil, err
	}

	sub, err := topic.topic.Subscribe()

	if err != nil {
		return nil, err
	}

	return &pubsubSubscription{
		s: sub,
	}, nil
}

func (p *PubSubManager) getOrJoinTopic(ctx context.Context, s string) (*pubsubTopic, error) {
	var err error

	p.m.Lock()
	defer p.m.Unlock()

	existing := p.topics[s]

	if existing != nil {
		return existing, nil
	}

	topic, err := p.ps.Join(s)

	if err != nil {
		return nil, err
	}

	existing = &pubsubTopic{
		topic: topic,
	}

	p.topics[s] = existing

	return existing, nil
}

type pubsubTopic struct {
	topic *pubsub.Topic
	refs  int
}

func (t *pubsubTopic) ref() {
	t.refs++
}

func (t *pubsubTopic) unref() {
	t.refs--
}

type pubsubSubscription struct {
	s *pubsub.Subscription
}

func (s *pubsubSubscription) Next(ctx context.Context) (iface.PubSubMessage, error) {
	m, err := s.s.Next(ctx)

	if err != nil {
		return nil, err
	}

	return pubusubMessage{m: m}, nil
}

func (s *pubsubSubscription) Close() error {
	s.s.Cancel()

	return nil
}

type pubusubMessage struct {
	m *pubsub.Message
}

func (p pubusubMessage) From() peer.ID {
	return p.m.ReceivedFrom
}

func (p pubusubMessage) Data() []byte {
	return p.m.Data
}

func (p pubusubMessage) Seq() []byte {
	return p.m.Seqno
}

func (p pubusubMessage) Topics() []string {
	if p.m.Topic == nil {
		return nil
	}

	return []string{*p.m.Topic}
}
