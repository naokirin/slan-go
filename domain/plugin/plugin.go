package plugin

import (
	"github.com/naokirin/slan-go/domain/slack"
)

// Config is plugin config for initialization
type Config struct {
	MentionName string
	Data        map[interface{}]interface{}
}

// CheckEnabledAdminUser checks message for admin user
func (c *Config) CheckEnabledAdminUser(msg slack.Message) bool {
	v, ok := c.Data["admin"]
	if !ok {
		return true
	}
	admin := v.(string)
	return admin == "*" || admin == msg.User
}

// CheckEnabledMessage checks message enabled
func (c *Config) CheckEnabledMessage(msg slack.Message) bool {
	return c.checkEnabledChannel(msg) && c.checkEnabledUser(msg)
}

func (c *Config) checkEnabledChannel(msg slack.Message) bool {
	v, ok := c.Data["channel"]
	if !ok {
		return true
	}
	channel := v.(string)
	return channel == "*" || channel == msg.Channel
}

func (c *Config) checkEnabledUser(msg slack.Message) bool {
	v, ok := c.Data["user"]
	if !ok {
		return true
	}
	user := v.(string)
	return user == "*" || user == msg.User
}
