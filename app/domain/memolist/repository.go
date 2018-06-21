package memolist

// Memo is memo interface
type Memo interface {
	GetID() uint
	GetUser() string
	GetText() string
	GetKind() string
}

// Repository is memolist repository interface
type Repository interface {
	All(kind string, user string) []Memo
	DeleteAll(kind string, user string)
	Add(kind string, user string, text string)
	Delete(m Memo) bool
}
