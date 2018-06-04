package memolist

import (
	"strconv"
	"strings"

	"github.com/naokirin/slan-go/domain/plugin"
	"github.com/naokirin/slan-go/domain/slack"
)

// GeneratePluginGoroutine generate memolist process
func GeneratePluginGoroutine(config plugin.Config, repository Repository, sendMessage func(string, string), in chan slack.Message) {
	plugin := &Plugin{
		MentionName: config.MentionName,
		SendMessage: sendMessage,
		Repository:  repository,
	}
	go func() {
		for msg := range in {
			if config.CheckEnabledMessage(msg) {
				plugin.ReceiveMessage(msg)
			}
		}
	}()
}

// Plugin is memolist plugin implementation
type Plugin struct {
	MentionName string
	SendMessage func(string, string)
	Repository  Repository
}

// ReceiveMessage processes memolist plugin for a received message
func (p *Plugin) ReceiveMessage(msg slack.Message) {
	if p.checkMessage(msg.Text, "list") {
		p.showList(msg)
	} else if p.checkMessage(msg.Text, "add") {
		p.addMemo(msg)
	} else if p.checkMessage(msg.Text, "delete") {
		p.deleteMemo(msg)
	}
}

func (p *Plugin) checkMessage(text string, subcommand string) bool {
	return strings.HasPrefix(text, "@"+p.MentionName+" memo."+subcommand)
}

func (p *Plugin) showList(msg slack.Message) {
	result := ""
	memolist := p.Repository.All(msg.User)
	for i, memo := range memolist {
		result = result + strconv.Itoa(i+1) + ". " + memo.GetText() + "\n"
		i++
	}
	if result == "" {
		result = "登録されたメモはありません"
	}
	p.SendMessage(result, msg.Channel)
}

func (p *Plugin) addMemo(msg slack.Message) {
	content := strings.SplitN(msg.Text, " ", 3)
	if len(content) >= 3 {
		contents := strings.Split(content[2], "\n")
		for _, c := range contents {
			p.Repository.Add(msg.User, c)
		}
		p.SendMessage("メモに追加しました", msg.Channel)
	} else {
		p.SendMessage("メモに追加できませんでした", msg.Channel)
	}
}

func (p *Plugin) deleteMemo(msg slack.Message) {
	all := p.Repository.All(msg.User)
	if len(all) <= 0 {
		p.SendMessage("登録されたメモがありません", msg.Channel)
	}
	memos := map[int]Memo{}
	for i, m := range all {
		memos[i] = m
	}

	content := strings.SplitN(msg.Text, " ", 3)
	if len(content) >= 3 {
		if content[2] == "all" {
			p.Repository.DeleteAll(msg.User)
			p.SendMessage("メモを削除しました", msg.Channel)
			return
		}
		result := false
		indexes := strings.Split(content[2], " ")
		for _, i := range indexes {
			index, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				p.SendMessage("メモ("+i+")を削除できませんでした", msg.Channel)
			} else {
				v, ok := memos[int(index-1)]
				if ok {
					p.Repository.Delete(msg.User, v)
					result = true
				} else {
					p.SendMessage("メモ("+i+")を削除できませんでした", msg.Channel)
				}
			}
		}
		if result {
			p.SendMessage("メモを削除しました", msg.Channel)
			return
		}
	}
	p.SendMessage("メモを削除できませんでした", msg.Channel)
}
