package yaml

import (
	"sync"
)

type repository struct {
	mentionName      string
	defaultResponses []string
	slackToken       string
	location         string
	language         string
	plugins          []interface{}
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
		dr, ok := data["default_responses"]
		defaultResponses := make([]string, 0)
		if ok {
			for _, d := range dr.([]interface{}) {
				defaultResponses = append(defaultResponses, d.(string))
			}
		}
		instance = &repository{
			mentionName:      data["mention_name"].(string),
			defaultResponses: defaultResponses,
			slackToken:       data["slack_token"].(string),
			plugins:          data["plugins"].([]interface{}),
			location:         data["location"].(string),
			language:         data["language"].(string),
		}
	})
	return ConfigurationRepository{*instance}
}

// GetMentionName returns the mention name
func (r *ConfigurationRepository) GetMentionName() string {
	return r.mentionName
}

// GetLanguage returns language
func (r *ConfigurationRepository) GetLanguage() string {
	return r.language
}

// GetDefaultResponses returns default responses
func (r *ConfigurationRepository) GetDefaultResponses() []string {
	return r.defaultResponses
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
