package slack

import (
	"log"
	"os"

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
func (client Client) SendMessage(msg string, channel string) s.Result {
	c, t, _, e := client.api.SendMessage(channel, sl.MsgOptionText(msg, true))
	return s.Result{
		Channel:   c,
		Timestamp: t,
		Err:       e,
	}
}

// UpdateMessage update message for existing one
func (client Client) UpdateMessage(text string, channel string, timestamp string) s.Result {
	c, t, _, e := client.api.SendMessage(channel, sl.MsgOptionUpdate(timestamp), sl.MsgOptionText(text, true))
	return s.Result{
		Channel:   c,
		Timestamp: t,
		Err:       e,
	}
}

func createAttachments(attachments []s.Attachment) []sl.Attachment {
	as := make([]sl.Attachment, 0, len(attachments))
	for _, at := range attachments {
		fields := make([]sl.AttachmentField, 0)
		for _, field := range at.Fields {
			fields = append(fields, sl.AttachmentField{
				Title: field.Title,
				Value: field.Value,
			})
		}

		a := sl.Attachment{
			Pretext: at.Pretext,
			Color:   at.Color,
			Fields:  fields,
		}
		as = append(as, a)
	}

	return as
}

// SendAttachments send attachment
func (client Client) SendAttachments(username string, attachments []s.Attachment, channel string) s.Result {
	params := sl.PostMessageParameters{
		Attachments: createAttachments(attachments),
		Username:    username,
		LinkNames:   1,
	}

	c, t, e := client.api.PostMessage(channel, "", params)
	return s.Result{
		Channel:   c,
		Timestamp: t,
		Err:       e,
	}
}

// UpdateAttachments update attachments for existing message
func (client Client) UpdateAttachments(username string, attachments []s.Attachment, channel string, timestamp string) s.Result {
	ats := createAttachments(attachments)
	attachmentsOpt := sl.MsgOptionAttachments(ats...)
	opts := []sl.MsgOption{
		attachmentsOpt,
		sl.MsgOptionUpdate(timestamp),
		sl.MsgOptionUser(username),
	}
	c, t, _, e := client.api.SendMessage(
		channel,
		opts...,
	)
	return s.Result{
		Channel:   c,
		Timestamp: t,
		Err:       e,
	}
}

// UploadFile upload file to slack
func (client Client) UploadFile(title string, path string, channel string) (string, error) {
	reader, err := os.Open(path)
	defer reader.Close()
	if err != nil {
		return "", err
	}
	params := sl.FileUploadParameters{
		Title:          title,
		InitialComment: "",
		Filename:       path,
		Reader:         reader,
		Channels:       []string{channel},
	}
	file, err := client.api.UploadFile(params)
	if err != nil {
		return "", err
	}
	return file.URL, err
}

// AddReaction add reaction to message
func (client Client) AddReaction(channel string, timestamp string, reaction string) error {
	msgRef := sl.NewRefToMessage(channel, timestamp)
	return client.api.AddReaction(reaction, msgRef)
}

// RemoveReaction remove reaction to message
func (client Client) RemoveReaction(channel string, timestamp string, reaction string) error {
	msgRef := sl.NewRefToMessage(channel, timestamp)
	return client.api.RemoveReaction(reaction, msgRef)
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
func (client Client) GenerateReceivedEventChannel() s.ReceivedChans {
	out := s.ReceivedChans{
		Message:         make(chan s.Message),
		ReactionAdded:   make(chan s.Reaction),
		ReactionRemoved: make(chan s.Reaction),
	}

	go func(chans s.ReceivedChans) {
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
				out.Message <- message
			case *sl.ReactionAddedEvent:
				// message と同じ情報を返せるようにする
				reaction := s.Reaction{
					Type:     ev.Type,
					User:     ev.User,
					ItemUser: ev.ItemUser,
					Item: s.ReactionItem{
						Type:        ev.Item.Type,
						Channel:     ev.Item.Channel,
						File:        ev.Item.File,
						FileComment: ev.Item.FileComment,
						Timestamp:   ev.Item.Timestamp,
					},
					Reaction:       ev.Reaction,
					EventTimestamp: ev.EventTimestamp,
				}
				out.ReactionAdded <- reaction
			case *sl.ReactionRemovedEvent:
				// message と同じ情報を返せるようにする
				reaction := s.Reaction{
					Type:     ev.Type,
					User:     ev.User,
					ItemUser: ev.ItemUser,
					Item: s.ReactionItem{
						Type:        ev.Item.Type,
						Channel:     ev.Item.Channel,
						File:        ev.Item.File,
						FileComment: ev.Item.FileComment,
						Timestamp:   ev.Item.Timestamp,
					},
					Reaction:       ev.Reaction,
					EventTimestamp: ev.EventTimestamp,
				}
				out.ReactionRemoved <- reaction
			default:
			}

		}
		close(out.Message)
		close(out.ReactionAdded)
		close(out.ReactionRemoved)
	}(out)

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

// GetBotUserID returns bot id
func (client Client) GetBotUserID() string {
	info := client.rtm.GetInfo()
	return info.User.ID
}

// GetEmoji returns all emoji
func (client Client) GetEmoji() map[string]string {
	result, err := client.api.GetEmoji()
	if err != nil {
		log.Printf("%v", err)
		return make(map[string]string)
	}
	return result
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

// GetUserName returns user display name
func (client Client) GetUserName(id string) (string, error) {
	user, err := client.rtm.GetUserInfo(id)
	if err != nil {
		return "", err
	}
	return user.Profile.DisplayName, nil
}
