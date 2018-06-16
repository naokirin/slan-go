package repository

// Repository is interface of configuration
type Repository interface {
	GetMentionName() string
	GetLanguage() string
	GetSlackToken() string
	GetLocation() string
	GetPlugins() map[string]interface{}
}
