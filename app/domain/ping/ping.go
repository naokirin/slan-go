package ping

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
)

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.Generator = (*Generator)(nil)

const defaultSubcommand = "ping"

var defaultResponses = []string{"pong"}

// Plugin is ping plugin implementation
type Plugin struct {
	mentionName string
	client      slack.Client
	responses   []string
	config      plugin.Config
}

// Generator is ping plugin generator
type Generator struct{}

// Generate generate memolist process
func (g *Generator) Generate(config plugin.Config, client slack.Client) plugin.Plugin {
	return &Plugin{
		mentionName: config.MentionName,
		client:      client,
		responses:   responses(config),
		config:      config,
	}
}

// ReceiveReactionAdded run received reaction_added
func (p *Plugin) ReceiveReactionAdded(reactionAdded slack.Reaction) {
}

// ReceiveReactionRemoved run received reaction_added
func (p *Plugin) ReceiveReactionRemoved(reactionRemoved slack.Reaction) {
}

// ReceiveMessage processes memolist plugin for a received message
func (p *Plugin) ReceiveMessage(msg slack.Message) bool {
	if !p.config.CheckEnabledMessage(msg) {
		return false
	}
	if p.checkMessage(msg.Text) {
		p.client.SendMessage(p.selectResponse(), msg.Channel)
		return true
	}
	return false
}

func (p *Plugin) checkMessage(text string) bool {
	subcommand := p.config.GetSubcommand(defaultSubcommand)
	return strings.HasPrefix(text, fmt.Sprintf("@%s %s", p.mentionName, subcommand))
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
	n := rand.Intn(len(p.responses))
	return p.responses[n]
}
