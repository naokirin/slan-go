package lunch

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/divan/num2words"
	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
)

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.Generator = (*Generator)(nil)

const defaultSubcommand = "lunch"

// Repository is interface for lunch data
type Repository interface {
	GetRows(sheetID string, readRange string, secretPath string, tokenPath string) [][]string
}

// Plugin for lunch choice
type Plugin struct {
	mentionName string
	client      slack.Client
	config      plugin.Config
	repository  Repository
}

// Generator for lunch choice plugin generation
type Generator struct {
	Repository Repository
}

// ReceiveMessage runs received message
func (p *Plugin) ReceiveMessage(msg slack.Message) bool {
	if !p.config.CheckEnabledMessage(msg) {
		return false
	}
	if p.checkMessage(msg.Text) {
		p.client.SendMessage(p.choiceLunches(), msg.Channel)
		return true
	}
	return false
}

func (p *Plugin) checkMessage(text string) bool {
	subcommand := p.config.GetSubcommand(defaultSubcommand)
	return strings.HasPrefix(text, fmt.Sprintf("@%s %s", p.mentionName, subcommand))
}

func (p *Plugin) choiceLunches() string {
	id, ok := p.config.Data["sheet_id"]
	if !ok {
		log.Printf("lunch plugin requires sheet_id setting.")
		return "...みつからなかったです!!"
	}

	ranges, ok := p.config.Data["ranges"]
	if !ok {
		log.Printf("lunch plugin requires ranges setting.")
		return "...みつからなかったです!!"
	}

	rs := make([][][]string, 0)
	for _, r := range ranges.([]interface{}) {
		rs = append(rs, p.repository.GetRows(id.(string), r.(string), p.getSecretPath(), p.getTokenPath()))
	}
	rand.Seed(time.Now().UnixNano())
	result := ""
	for i, r := range rs {
		l := len(r)
		if l == 0 {
			continue
		}
		n := rand.Intn(l)
		result += fmt.Sprintf(":%s: %s\n", num2words.Convert(i+1), r[n][0])
	}
	if len(result) == 0 {
		return "...みつからなかったです!!"
	}
	return "ランチ候補です!!\n" + result
}

func (p *Plugin) getTokenPath() string {
	path, ok := p.config.Data["token_file"]
	if !ok {
		return ""
	}
	return path.(string)
}

func (p *Plugin) getSecretPath() string {
	path, ok := p.config.Data["secret_file"]
	if !ok {
		return ""
	}
	return path.(string)
}

// Generate generates lunch choice plugin
func (g *Generator) Generate(config plugin.Config, client slack.Client) plugin.Plugin {
	return &Plugin{
		mentionName: config.MentionName,
		client:      client,
		config:      config,
		repository:  g.Repository,
	}
}
