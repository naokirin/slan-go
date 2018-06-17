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
	SendAttachments(name string, attachments []Attachment, channel string)
	UploadFile(title string, path string, channel string) (string, error)
	GenerateReceivedEventChannel() chan Message
	GetBotName() string
	GetEmoji() map[string]string
	ConvertChannelNameToID(name string) (string, bool)
}
