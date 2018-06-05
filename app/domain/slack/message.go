package slack

// Message is a received message type
type Message struct {
	Type      string
	User      string
	Text      string
	TimeStamp string
	Channel   string
}
