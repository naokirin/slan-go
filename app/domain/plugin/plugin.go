package plugin

import (
	"github.com/naokirin/slan-go/app/domain/slack"
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
	for _, admin := range v.([]interface{}) {
		a := admin.(string)
		return a == "*" || a == msg.User
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
		if c == "*" || c == msg.Channel {
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
		if u == "*" || u == msg.User {
			return true
		}
	}
	return false
}
