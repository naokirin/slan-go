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

// ReceiveMessage processes memolist plugin for a received message
func (p *Plugin) ReceiveMessage(msg slack.Message) {
	if !p.config.CheckEnabledMessage(msg) {
		return
	}
	command := p.config.GetSubcommand("memo")
	if p.checkMessage(msg.Text, command, "list") {
		p.showList(msg)
	} else if p.checkMessage(msg.Text, command, "add") {
		p.addMemo(msg)
	} else if p.checkMessage(msg.Text, command, "delete") {
		p.deleteMemo(msg)
	}
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
		result = "メモはありません!!"
	}
	p.client.SendMessage(result, msg.Channel)
}

func (p *Plugin) addMemo(msg slack.Message) {
	content := strings.SplitN(msg.Text, " ", 3)
	if len(content) >= 3 {
		contents := strings.Split(content[2], "\n")
		for _, c := range contents {
			p.repository.Add(p.kind, msg.User, c)
		}
		p.client.SendMessage("メモしました!", msg.Channel)
	} else {
		p.client.SendMessage("メモできませんでした...", msg.Channel)
	}
}

func (p *Plugin) deleteMemo(msg slack.Message) {
	all := p.repository.All(p.kind, msg.User)
	if len(all) <= 0 {
		p.client.SendMessage("メモはないです!", msg.Channel)
	}
	memos := map[int]Memo{}
	for i, m := range all {
		memos[i] = m
	}

	content := strings.SplitN(msg.Text, " ", 3)
	if len(content) >= 3 {
		if content[2] == "all" {
			p.repository.DeleteAll(p.kind, msg.User)
			p.client.SendMessage("メモを削除しました", msg.Channel)
			return
		}
		result := false
		indexes := strings.Split(content[2], " ")
		for _, i := range indexes {
			index, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				p.client.SendMessage("メモ("+i+")を削除できませんでした", msg.Channel)
			} else {
				v, ok := memos[int(index-1)]
				if ok {
					p.repository.Delete(v)
					result = true
				} else {
					p.client.SendMessage("メモ("+i+")を削除できませんでした", msg.Channel)
				}
			}
		}
		if result {
			p.client.SendMessage("メモを削除しました", msg.Channel)
			return
		}
	}
	p.client.SendMessage("メモを削除できませんでした", msg.Channel)
}
