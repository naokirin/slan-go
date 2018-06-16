package lgtmize

import (
	"fmt"
	"strconv"
	"strings"
	"os"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
)

// LGTMize is interface for lgtm creation
type LGTMize interface {
	CreateLGTM(url string, size int, color string) (string, error)
}

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.Generator = (*Generator)(nil)

const defaultSubcommand = "lgtmize"

// Plugin is lgtmize plugin implementation
type Plugin struct {
	mentionName string
	client      slack.Client
	config      plugin.Config
	lgtmize     LGTMize
}

// Generator is lgtmize plugin generator
type Generator struct {
	LGTMize LGTMize
}

// Generate generate lgtmize process
func (g *Generator) Generate(config plugin.Config, client slack.Client) plugin.Plugin {
	return &Plugin{
		mentionName: config.MentionName,
		client:      client,
		config:      config,
		lgtmize:     g.LGTMize,
	}
}

// ReceiveMessage processes lgtm plugin for a received message
func (p *Plugin) ReceiveMessage(msg slack.Message) bool {
	if !p.config.CheckEnabledMessage(msg) {
		return false
	}
	if p.checkMessage(msg.Text) {
		ss := strings.SplitN(msg.Text, " ", 5)
		if len(ss) < 3 {
			message := p.config.ResponseTemplates.GetText("error_message_for_generate_image", nil)
			p.client.SendMessage(message, msg.Channel)
			return true
		}
		if len(ss[2]) < 3 {
			message := p.config.ResponseTemplates.GetText("error_message_for_generate_image", nil)
			p.client.SendMessage(message, msg.Channel)
			return true
		}

		color := "black"
		if len(ss) >= 4 {
			color = ss[3]
		}
		size := 500
		if len(ss) >= 5 {
			var err error
			size, err = strconv.Atoi(ss[4])
			if err != nil {
				size = 500
			}
		}
		go func(url string, color string, channel string) {
			path, err := p.lgtmize.CreateLGTM(url, size, color)
			if err != nil {
				message := p.config.ResponseTemplates.GetText("error_message_for_generate_image", nil)
				p.client.SendMessage(message, channel)
				return
			}
			defer os.Remove(path)
			imageURL, err := p.client.UploadFile("LGTM", path, channel)
			if err != nil {
				message := p.config.ResponseTemplates.GetText("error_message_for_upload", nil)
				p.client.SendMessage(message, channel)
				return
			}
			p.client.SendMessage(imageURL, channel)
		}(ss[2][1:len(ss[2])-1], color, msg.Channel)
		return true
	}
	return false
}

func (p *Plugin) checkMessage(text string) bool {
	subcommand := p.config.GetSubcommand(defaultSubcommand)
	return strings.HasPrefix(text, fmt.Sprintf("@%s %s", p.mentionName, subcommand))
}
