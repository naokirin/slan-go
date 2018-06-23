package memolist

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
)

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.Generator = (*Generator)(nil)

// Plugin is memolist plugin implementation
type Plugin struct {
	mentionName string
	client      slack.Client
	repository  Repository
	kind        string
	config      plugin.Config
}

// Generator is memolist plugin generator
type Generator struct {
	Repository Repository
}

// Generate generate memolist process
func (g *Generator) Generate(config plugin.Config, client slack.Client) plugin.Plugin {
	k, ok := config.Data["kind"]
	kind := ""
	if ok {
		kind = k.(string)
	}
	return &Plugin{
		mentionName: config.MentionName,
		client:      client,
		kind:        kind,
		repository:  g.Repository,
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
	command := p.config.GetSubcommand("memo")
	if p.checkMessage(msg.Text, command, "list") {
		p.showList(msg)
		return true
	}
	if p.checkMessage(msg.Text, command, "add") {
		p.addMemo(msg)
		return true
	}
	if p.checkMessage(msg.Text, command, "delete") {
		p.deleteMemo(msg)
		return true
	}
	return false
}

func (p *Plugin) checkMessage(text string, command string, subcommand string) bool {
	return strings.HasPrefix(text, fmt.Sprintf("@%s %s.%s", p.mentionName, command, subcommand))
}

func (p *Plugin) showList(msg slack.Message) {
	result := ""
	memolist := p.repository.All(p.kind, msg.User)
	for i, memo := range memolist {
		result = result + strconv.Itoa(i+1) + ". " + memo.GetText() + "\n"
		i++
	}
	if result == "" {
		result = p.config.ResponseTemplates.GetText("not_found_memo", nil)
	}
	p.client.SendMessage(result, msg.Channel)
}

func (p *Plugin) addMemo(msg slack.Message) {
	content := strings.SplitN(msg.Text, " ", 3)
	message := ""
	if len(content) >= 3 {
		contents := strings.Split(content[2], "\n")
		for _, c := range contents {
			p.repository.Add(p.kind, msg.User, c)
		}
		message = p.config.ResponseTemplates.GetText("add_memo", nil)
	} else {
		message = p.config.ResponseTemplates.GetText("could_not_add_memo", nil)
	}
	p.client.SendMessage(message, msg.Channel)
}

func (p *Plugin) deleteMemo(msg slack.Message) {
	all := p.repository.All(p.kind, msg.User)
	if len(all) <= 0 {
		message := p.config.ResponseTemplates.GetText("not_found_memo_when_delete", nil)
		p.client.SendMessage(message, msg.Channel)
		return
	}
	memos := map[int]Memo{}
	for i, m := range all {
		memos[i] = m
	}

	content := strings.SplitN(msg.Text, " ", 3)
	if len(content) >= 3 {
		if content[2] == "all" {
			p.repository.DeleteAll(p.kind, msg.User)
			message := p.config.ResponseTemplates.GetText("delete_memo", nil)
			p.client.SendMessage(message, msg.Channel)
			return
		}
		result := false
		indexes := strings.Split(content[2], " ")
		for _, i := range indexes {
			index, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				message := p.config.ResponseTemplates.GetText("could_not_delete_specified_delete", map[string]string{"Number": i})
				p.client.SendMessage(message, msg.Channel)
			} else {
				v, ok := memos[int(index-1)]
				if ok {
					p.repository.Delete(v)
					result = true
				} else {
					message := p.config.ResponseTemplates.GetText("could_not_delete_specified_delete", map[string]string{"Number": i})
					p.client.SendMessage(message, msg.Channel)
				}
			}
		}
		if result {
			message := p.config.ResponseTemplates.GetText("delete_memo", nil)
			p.client.SendMessage(message, msg.Channel)
			return
		}
	}
	message := p.config.ResponseTemplates.GetText("could_not_delete_memo", nil)
	p.client.SendMessage(message, msg.Channel)
}
