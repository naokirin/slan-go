package plugin

import (
	"github.com/naokirin/slan-go/app/domain/slack"
)

// Plugin is interface of received message plugin
type Plugin interface {
	ReceiveMessage(msg slack.Message) bool
}

// Generator is interface of generate plugin
type Generator interface {
	Generate(Config, slack.Client) Plugin
}

// Config is plugin config for initialization
type Config struct {
	MentionName string
	Data        map[interface{}]interface{}
}

// GetSubcommand returns subcommand name
func (c *Config) GetSubcommand(defaultName string) string {
	v, ok := c.Data["subcommand"]
	if !ok {
		return defaultName
	}
	return v.(string)
}

// CheckEnabledAdminUser checks message for admin user
func (c *Config) CheckEnabledAdminUser(msg slack.Message) bool {
	v, ok := c.Data["admin"]
	if !ok {
		return true
	}
	for _, admin := range v.([]interface{}) {
		a := admin.(string)
		return a == "*" || a == msg.UserName
	}
	return false
}

// CheckEnabledMessage checks message enabled
func (c *Config) CheckEnabledMessage(msg slack.Message) bool {
	return c.checkEnabledChannel(msg) && c.checkEnabledUser(msg)
}

func (c *Config) checkEnabledChannel(msg slack.Message) bool {
	v, ok := c.Data["channels"]
	if !ok {
		return true
	}
	for _, channel := range v.([]interface{}) {
		c := channel.(string)
		if c == "*" || c == msg.ChannelName {
			return true
		}
	}
	return false
}

func (c *Config) checkEnabledUser(msg slack.Message) bool {
	v, ok := c.Data["users"]
	if !ok {
		return true
	}
	for _, user := range v.([]interface{}) {
		u := user.(string)
		if u == "*" || u == msg.UserName {
			return true
		}
	}
	return false
}
