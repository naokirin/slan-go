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
				botUser := client.rtm.GetInfo().User
				channel, err := client.rtm.GetChannelInfo(ev.Channel)
				channelName := ev.Username
				if err == nil {
					channelName = channel.Name
				}
				message := s.Message{
					Type:        ev.Type,
					BotName:     botUser.Name,
					User:        ev.User,
					UserName:    ev.Username,
					Text:        ev.Text,
					TimeStamp:   ev.Timestamp,
					Channel:     ev.Channel,
					ChannelName: channelName,
				}
				out <- message
			default:
			}

		}
		close(out)
	}()

	return out
}
