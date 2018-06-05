package memolist

// Memo is memo interface
type Memo interface {
	GetUser() string
	GetText() string
}

// Repository is memolist repository interface
type Repository interface {
	All(user string) []Memo
	DeleteAll(user string)
	Add(user string, text string)
	Delete(user string, m Memo) bool
}
