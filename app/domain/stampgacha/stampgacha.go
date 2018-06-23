package stampgacha

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
)

// Repository is interface for stamp gacha
type Repository interface {
	LastLottingTime(user string) (time.Time, bool)
	SaveLottingTime(user string)
}

// ConfigRepository is interface for stamp gacha emoji list
type ConfigRepository interface {
	GetEmojiListRepository(path string) ConfigRepository
	GetEmojiList() []string
}

const defaultSubcommand = "stamp_gacha"

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.Generator = (*Generator)(nil)

// Plugin is stampgacha plugin implementation
type Plugin struct {
	mentionName      string
	client           slack.Client
	repository       Repository
	configRepository ConfigRepository
	config           plugin.Config
	stamps           []string
}

// Generator is stampgacha plugin generator
type Generator struct {
	Repository       Repository
	ConfigRepository ConfigRepository
}

// Generate generate stampgacha process
func (g *Generator) Generate(config plugin.Config, client slack.Client) plugin.Plugin {
	ss, ok := config.Data["stamps"]
	crep := make([]string, 0)
	if !ok {
		log.Println("stamp_gacha requires 'stamps'")
	} else {
		crep = g.ConfigRepository.GetEmojiListRepository(ss.(string)).GetEmojiList()
	}
	return &Plugin{
		mentionName: config.MentionName,
		client:      client,
		repository:  g.Repository,
		config:      config,
		stamps:      crep,
	}
}

// ReceiveReactionAdded run received reaction_added
func (p *Plugin) ReceiveReactionAdded(reactionAdded slack.Reaction) {
}

// ReceiveReactionRemoved run received reaction_added
func (p *Plugin) ReceiveReactionRemoved(reactionRemoved slack.Reaction) {
}

// ReceiveMessage processes stampgacha plugin for a received message
func (p *Plugin) ReceiveMessage(msg slack.Message) bool {
	if !p.config.CheckEnabledMessage(msg) {
		return false
	}
	if p.checkMessage(msg.Text) {
		p.draw(msg)
		return true
	}
	return false
}

func (p *Plugin) checkMessage(text string) bool {
	subcommand := p.config.GetSubcommand(defaultSubcommand)
	return strings.HasPrefix(text, fmt.Sprintf("@%s %s", p.mentionName, subcommand))
}

func (p *Plugin) draw(msg slack.Message) {
	stampLen := len(p.stamps)
	if stampLen == 0 {
		message := p.config.ResponseTemplates.GetText("not_found_stamps", nil)
		p.client.SendMessage(message, msg.Channel)
		return
	}
	lastLottingTime, ok := p.repository.LastLottingTime(msg.User)
	if ok {
		diff := time.Since(lastLottingTime)
		if diff < time.Minute*5 {
			m := map[string]string{"Duration": "5"}
			message := p.config.ResponseTemplates.GetText("not_allowed_for_duration", m)
			p.client.SendMessage(message, msg.Channel)
			return
		}
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(stampLen)
	m := map[string]string{"Stamp": p.stamps[n]}
	message := p.config.ResponseTemplates.GetText("show_stamp", m)
	p.client.SendMessage(message, msg.Channel)
	p.repository.SaveLottingTime(msg.User)
}
