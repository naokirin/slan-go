package slack

// Client is interface slack client using plugin
type Client interface {
	SendMessage(msg string, channel string)
	GenerateReceivedEventChannel() chan Message
}
