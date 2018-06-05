package ping

import (
	"math/rand"
	"strings"
	"time"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
)

var defaultResponses = []string{"pong"}

// Plugin is ping plugin implementation
type Plugin struct {
	MentionName string
	SendMessage func(string, string)
	Responses   []string
}

// GeneratePluginGoroutine generate memolist process
func GeneratePluginGoroutine(config plugin.Config, sendMessage func(string, string), in chan slack.Message) {
	plugin := &Plugin{
		MentionName: config.MentionName,
		SendMessage: sendMessage,
		Responses:   responses(config),
	}

	go func() {
		for msg := range in {
			if config.CheckEnabledMessage(msg) {
				plugin.ReceiveMessage(msg)
			}
		}
	}()
}

// ReceiveMessage processes memolist plugin for a received message
func (p *Plugin) ReceiveMessage(msg slack.Message) {
	if p.checkMessage(msg.Text) {
		p.SendMessage(p.selectResponse(), msg.Channel)
	}
}

func (p *Plugin) checkMessage(text string) bool {
	return strings.HasPrefix(text, "@"+p.MentionName+" ping")
}

func responses(config plugin.Config) []string {
	list, ok := config.Data["responses"]
	if !ok {
		return defaultResponses
	}
	responses := list.([]interface{})
	if len(responses) <= 0 {
		return defaultResponses
	}
	result := []string{}
	for _, r := range responses {
		result = append(result, r.(string))
	}
	return result
}

func (p *Plugin) selectResponse() string {

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(p.Responses))
	return p.Responses[n]
}
