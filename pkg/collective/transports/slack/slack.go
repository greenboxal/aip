package slack

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/slack-go/slack"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/pkg/collective"
)

type Transport struct {
	m sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	api *slack.Client
	rtm *slack.RTM

	channelNameCache map[string]string
	userNameCache    map[string]string

	incoming chan collective.Message
}

func NewTransport(lc fx.Lifecycle) *Transport {
	ctx, cancel := context.WithCancel(context.Background())

	t := &Transport{
		ctx:    ctx,
		cancel: cancel,

		channelNameCache: map[string]string{},
		userNameCache:    map[string]string{},

		incoming: make(chan collective.Message, 16),
	}

	botToken := os.Getenv("SLACK_BOT_USER_TOKEN")

	t.api = slack.New(botToken)
	t.rtm = t.api.NewRTM()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, _, err := t.rtm.ConnectRTMContext(ctx)

			if err != nil {
				return err
			}

			go func() {
				t.rtm.ManageConnection()
			}()

			go func() {
				for ev := range t.rtm.IncomingEvents {
					switch evt := ev.Data.(type) {
					case *slack.ConnectingEvent:
						fmt.Println("Connecting to Slack with Socket Mode...")
					case *slack.ConnectionErrorEvent:
						fmt.Println("Connection failed. Retrying later...")
					case *slack.ConnectedEvent:
						fmt.Println("Connected to Slack with Socket Mode.")

					case *slack.MessageEvent:
						t.incoming <- collective.Message{
							ID:      evt.ClientMsgID,
							Channel: t.resolveChannel(evt),
							From:    t.resolveUser(evt),
							Text:    evt.Text,
						}
					}
				}
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			return t.Close()
		},
	})

	return t
}

func (t *Transport) Subscribe(channel string) error {
	return nil
}

func (t *Transport) Incoming() <-chan collective.Message {
	return t.incoming
}

func (t *Transport) RouteMessage(ctx context.Context, msg collective.Message) error {
	_, _, err := t.rtm.PostMessage(
		msg.Channel,
		slack.MsgOptionText(msg.Text, true),
		slack.MsgOptionUsername(msg.From),
	)

	return err
}

func (t *Transport) Close() error {
	t.cancel()

	_ = t.rtm.Disconnect()

	return nil
}

func (t *Transport) resolveChannel(evt *slack.MessageEvent) string {
	if existing, ok := t.channelNameCache[evt.Channel]; ok {
		return existing
	}

	t.m.Lock()
	defer t.m.Unlock()

	ch, err := t.rtm.GetConversationInfo(&slack.GetConversationInfoInput{ChannelID: evt.Channel})

	if err != nil {
		return evt.Channel
	}

	t.channelNameCache[evt.Channel] = ch.Name

	return ch.Name
}

func (t *Transport) resolveUser(evt *slack.MessageEvent) string {
	if existing, ok := t.userNameCache[evt.User]; ok {
		return existing
	}

	t.m.Lock()
	defer t.m.Unlock()

	u, err := t.rtm.GetUserInfo(evt.User)

	if err != nil {
		return evt.Username
	}

	t.userNameCache[evt.User] = u.Name

	return u.Name
}
