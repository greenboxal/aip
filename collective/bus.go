package collective

import "sync"

// Channel is a channel that can be used to send messages to a single receiver.
type Channel interface {
	Send(msg Message) (Message, error)
}

// Bus is a channel that can be used to send messages to multiple receivers.
type Bus interface {
	Channel
}

// Message is a message that can be sent over a channel.
type Message struct {
	Message string
}

type ChannelDialer interface {
	Dial(name string) (Channel, error)
}

type Collective interface {
	Dialer() ChannelDialer

	Send(channel string, msg Message) (Message, error)
}

type basicCollective struct {
	m             sync.Mutex
	dialer        ChannelDialer
	knownChannels map[string]Channel
}

func NewBasicCollective(dialer ChannelDialer) (Collective, error) {
	return &basicCollective{
		dialer:        dialer,
		knownChannels: make(map[string]Channel),
	}, nil
}

func (b *basicCollective) Dialer() ChannelDialer {
	return b.dialer
}

func (b *basicCollective) Send(channel string, msg Message) (Message, error) {
	ch := b.knownChannels[channel]

	if ch == nil {
		var err error

		b.m.Lock()
		defer b.m.Unlock()

		ch, err = b.dialer.Dial(channel)

		if err != nil {
			return Message{}, err
		}

		b.knownChannels[channel] = ch
	}

	return ch.Send(msg)
}
