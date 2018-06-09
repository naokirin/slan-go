package slack

import (
	s "github.com/naokirin/slan-go/app/domain/slack"
	sl "github.com/nlopes/slack"
)

var _ s.Client = Client{}

// Client is a slack client wrapper
type Client struct {
	api *sl.Client
	rtm *sl.RTM
}

// SendMessage send message with RTM
func (client Client) SendMessage(msg string, channel string) {
	client.rtm.SendMessage(client.rtm.NewOutgoingMessage(msg, channel))
}

// SendAttachment send attachment
func (client Client) SendAttachment(name string, attachment s.Attachment, channel string) {

	fields := make([]sl.AttachmentField, 0)
	for _, field := range attachment.Fields {
		fields = append(fields, sl.AttachmentField{
			Title: field.Title,
			Value: field.Value,
		})
	}

	a := sl.Attachment{
		Pretext: attachment.Pretext,
		Color:   attachment.Color,
		Fields:  fields,
	}

	params := sl.PostMessageParameters{
		Attachments: []sl.Attachment{a},
		Username:    name,
	}

	client.api.PostMessage(channel, "", params)
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
				user, err := client.rtm.GetUserInfo(ev.User)
				userName := ev.Username
				if err == nil {
					userName = user.Profile.DisplayName
				}
				channel, err := client.rtm.GetChannelInfo(ev.Channel)
				channelName := userName
				if err == nil {
					channelName = channel.Name
				}
				message := s.Message{
					Type:        ev.Type,
					BotName:     botUser.Name,
					User:        ev.User,
					UserName:    userName,
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

// GetBotName returns bot display name
func (client Client) GetBotName() string {
	info := client.rtm.GetInfo()
	user, err := client.api.GetUserInfo(info.User.ID)
	if err != nil {
		return info.User.Name
	}
	return user.RealName
}

// ConvertChannelNameToID returns corresponding channel id
func (client Client) ConvertChannelNameToID(name string) (string, bool) {
	cs, err := client.api.GetChannels(false)
	if err != nil {
		return "", false
	}
	for _, c := range cs {
		if c.Name == name {
			return c.ID, true
		}
	}
	return "", false
}
