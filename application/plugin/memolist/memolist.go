package memolist

import (
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	// Register some standard stuff
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/naokirin/slan-go/domain/plugin"
	dslack "github.com/naokirin/slan-go/domain/slack"
	"github.com/naokirin/slan-go/infrastructure/slack"
)

type memo struct {
	gorm.Model
	User string
	Text string
}

// Plugin is memolist plugin implementation
type Plugin struct {
	mentionName string
}

// GeneratePluginGoroutine generate memolist process
func GeneratePluginGoroutine(config plugin.Config, client *slack.Client, in chan dslack.Message) {
	plugin := &Plugin{config.MentionName}
	go func() {
		for msg := range in {
			if config.CheckEnabledMessage(msg) {
				plugin.ReceiveMessage(client, msg)
			}
		}
	}()
}

// ReceiveMessage processes memolist plugin for a received message
func (p *Plugin) ReceiveMessage(client *slack.Client, msg dslack.Message) {
	if p.checkMessage(msg.Text, "list") {
		showList(client, msg)
	} else if p.checkMessage(msg.Text, "add") {
		addMemo(client, msg)
	} else if p.checkMessage(msg.Text, "delete") {
		deleteMemo(client, msg)
	}
}

func (p *Plugin) checkMessage(text string, subcommand string) bool {
	return strings.HasPrefix(text, "@"+p.mentionName+" memo."+subcommand)
}

func showList(client *slack.Client, msg dslack.Message) {
	result := ""
	memolist := all(msg.User)
	for i, memo := range memolist {
		result = result + strconv.Itoa(i+1) + ". " + memo.Text + "\n"
		i++
	}
	if result == "" {
		result = "登録されたメモはありません"
	}
	client.SendMessage(result, msg.Channel)
}

func addMemo(client *slack.Client, msg dslack.Message) {
	content := strings.SplitN(msg.Text, " ", 3)
	if len(content) >= 3 {
		contents := strings.Split(content[2], "\n")
		for _, c := range contents {
			add(msg.User, c)
		}
		client.SendMessage("メモに追加しました", msg.Channel)
	} else {
		client.SendMessage("メモに追加できませんでした", msg.Channel)
	}
}

func deleteMemo(client *slack.Client, msg dslack.Message) {
	all := all(msg.User)
	if len(all) <= 0 {
		client.SendMessage("登録されたメモがありません", msg.Channel)
	}
	memos := map[int]memo{}
	for i, m := range all {
		memos[i] = m
	}

	content := strings.SplitN(msg.Text, " ", 3)
	if len(content) >= 3 {
		if content[2] == "all" {
			deleteAll(msg.User)
			client.SendMessage("メモを削除しました", msg.Channel)
			return
		}
		result := false
		indexes := strings.Split(content[2], " ")
		for _, i := range indexes {
			index, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				client.SendMessage("メモ("+i+")を削除できませんでした", msg.Channel)
			} else {
				v, ok := memos[int(index-1)]
				if ok {
					delete(msg.User, v)
					result = true
				} else {
					client.SendMessage("メモ("+i+")を削除できませんでした", msg.Channel)
				}
			}
		}
		if result {
			client.SendMessage("メモを削除しました", msg.Channel)
			return
		}
	}
	client.SendMessage("メモを削除できませんでした", msg.Channel)
}
