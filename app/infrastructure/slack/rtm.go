package slack

import (
	s "github.com/naokirin/slan-go/app/domain/slack"
	sl "github.com/nlopes/slack"
)

// Client is a slack client wrapper
type Client struct {
	api *sl.Client
	rtm *sl.RTM
}

// SendMessage send message with RTM
func (client Client) SendMessage(msg string, channel string) {
	client.rtm.SendMessage(client.rtm.NewOutgoingMessage(msg, channel))
}

// CreateClient create slack client for RTM
func CreateClient(slackToken string) *Client {
	api := sl.New(slackToken)
	rtm := api.NewRTM()
	client := Client{api, rtm}
	go rtm.ManageConnection()
	return &client
}

// GenerateReceivedEventChannel generate a channel of received incomming event
func (client Client) GenerateReceivedEventChannel() chan s.Message {
	out := make(chan s.Message)
	go func() {
		for msg := range client.rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *sl.MessageEvent:
				var message s.Message
				message.Type = ev.Type
				message.User = ev.User
				message.Text = ev.Text
				message.TimeStamp = ev.Timestamp
				message.Channel = ev.Channel
				out <- message
			default:
			}

		}
		close(out)
	}()

	return out
}
