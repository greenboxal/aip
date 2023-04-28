package slack

import (
	"context"
	"fmt"
	"os"
	"regexp"
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

var messageHeaderRegex = regexp.MustCompile(":thread: (?P<thread_id>[^ ]+) \\[(?P<reply_to_id>[^]]+)]: (?P<text>.*)")

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
			return t.Start(ctx)
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
	var options []slack.MsgOption

	options = append(
		options,
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionUsername(msg.From),
		slack.MsgOptionMetadata(slack.SlackMetadata{
			EventType: "aip_say",
			EventPayload: map[string]interface{}{
				"mid":         msg.ID,
				"reply_to_id": msg.ReplyToID,
				"thread_id":   msg.ThreadID,
				"from":        msg.From,
			},
		}),
	)

	if msg.ThreadID != "" {
		options = append(options, slack.MsgOptionTS(msg.ThreadID))
	}

	_, _, err := t.rtm.PostMessage(msg.Channel, options...)

	return err
}

func (t *Transport) Close() error {
	t.cancel()

	_ = t.rtm.Disconnect()

	return nil
}

func (t *Transport) resolveChannel(evt *slack.Message) string {
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

func (t *Transport) resolveUser(evt *slack.Message) string {
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

func (t *Transport) Start(ctx context.Context) error {
	go func() {
		_, _, err := t.rtm.ConnectRTMContext(ctx)

		if err != nil {
			panic(err)
		}

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
				slackMsg := (*slack.Message)(evt)

				msg := collective.Message{
					ID:       slackMsg.Timestamp,
					ThreadID: slackMsg.ThreadTimestamp,
					Channel:  t.resolveChannel(slackMsg),
					From:     t.resolveUser(slackMsg),
					Text:     evt.Text,
				}

				if slackMsg.Metadata.EventType == "aip_say" {
					if slackMsg.Metadata.EventPayload != nil {
						msg.ReplyToID = slackMsg.Metadata.EventPayload["reply_to_id"].(string)
						msg.ThreadID = slackMsg.Metadata.EventPayload["thread_id"].(string)
						msg.ID = slackMsg.Metadata.EventPayload["id"].(string)
						msg.From = slackMsg.Metadata.EventPayload["from"].(string)
					}

					groups := messageHeaderRegex.FindStringSubmatch(slackMsg.Text)

					if groups != nil {
						msg.ThreadID = groups[1]
						msg.ReplyToID = groups[2]
						msg.Text = groups[3]
					}
				}

				t.incoming <- msg
			}
		}
	}()

	return nil
}
