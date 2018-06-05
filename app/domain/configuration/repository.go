package repository

// Repository is interface of configuration
type Repository interface {
	GetMentionName() string
	GetSlackToken() string
	GetPlugins() map[string]interface{}
}
