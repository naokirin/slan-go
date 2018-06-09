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

// Client is interface slack client using plugin
type Client interface {
	SendMessage(msg string, channel string)
	SendAttachment(name string, attachment Attachment, channel string)
	GenerateReceivedEventChannel() chan Message
	GetBotName() string
	ConvertChannelNameToID(name string) (string, bool)
}
