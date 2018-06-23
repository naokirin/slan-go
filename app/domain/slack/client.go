package slack

// AttachmentField for attachment
type AttachmentField struct {
	Title string
	Value string
}

// Attachment for attachment
type Attachment struct {
	Pretext string
	Color   string
	Fields  []AttachmentField
}

// ReceivedChans is slack reveived event chans
type ReceivedChans struct {
	Message         chan Message
	ReactionAdded   chan Reaction
	ReactionRemoved chan Reaction
}

// Result for sending message result
type Result struct {
	Channel   string
	Timestamp string
	Err       error
}

// Client is interface slack client using plugin
type Client interface {
	SendMessage(msg string, channel string) Result
	SendAttachments(username string, attachments []Attachment, channel string) Result
	UpdateMessage(text string, channel string, timestamp string) Result
	UpdateAttachments(username string, attachments []Attachment, channel string, timestamp string) Result
	UploadFile(title string, path string, channel string) (string, error)
	AddReaction(channel string, timestamp string, reaction string) error
	RemoveReaction(channel string, timestamp string, reaction string) error
	GenerateReceivedEventChannel() ReceivedChans
	GetBotName() string
	GetBotUserID() string
	GetEmoji() map[string]string
	GetUserName(id string) (string, error)
	ConvertChannelNameToID(name string) (string, bool)
}
