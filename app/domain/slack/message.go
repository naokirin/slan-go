package slack

// Message is a received message type
type Message struct {
	Type      string
	BotName   string
	User      string
	UserName  string
	Text      string
	TimeStamp string
	Channel   string
	ChannelName string
}
