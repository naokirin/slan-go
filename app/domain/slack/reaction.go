package slack

// ReactionItem is item for reaction_added or reaction_removed
type ReactionItem struct {
	Type        string
	Channel     string
	File        string
	FileComment string
	Timestamp   string
}

// Reaction is type for reaction_added or reaction_removed
type Reaction struct {
	Type           string
	User           string
	ItemUser       string
	Item           ReactionItem
	Reaction       string
	EventTimestamp string
}
