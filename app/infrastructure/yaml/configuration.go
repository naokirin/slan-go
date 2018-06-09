package yaml

import (
	"sync"
)

type repository struct {
	mentionName string
	slackToken  string
	location    string
	plugins     []interface{}
}

// ConfigurationRepository is configuration data
type ConfigurationRepository struct {
	repository
}

const confPath = "config/slan-go.conf"

var instance *repository
var once sync.Once

// GetConfigurationRepository returns configuration
func GetConfigurationRepository() ConfigurationRepository {
	once.Do(func() {
		data, err := ParseFromFile(confPath)
		if err != nil {
			return
		}
		instance = &repository{
			mentionName: data["mention_name"].(string),
			slackToken:  data["slack_token"].(string),
			plugins:     data["plugins"].([]interface{}),
			location:    data["location"].(string),
		}
	})
	return ConfigurationRepository{*instance}
}

// GetMentionName returns the mention name
func (r *ConfigurationRepository) GetMentionName() string {
	return r.mentionName
}

// GetSlackToken returns slack token from configuration
func (r *ConfigurationRepository) GetSlackToken() string {
	return r.slackToken
}

// GetLocation returns location name for scheduler
func (r *ConfigurationRepository) GetLocation() string {
	return r.location
}

// GetPlugins returns plugin configurations
func (r *ConfigurationRepository) GetPlugins() []interface{} {
	return r.plugins
}
