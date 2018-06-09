package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/schedule"
	"github.com/naokirin/slan-go/app/domain/slack"
)

var _ plugin.Plugin = (*Plugin)(nil)
var _ plugin.Generator = (*Generator)(nil)

// Plugin for calendar
type Plugin struct {
	config   plugin.Config
	client   slack.Client
	calendar Calendar
}

// ReceiveMessage run received message
func (p *Plugin) ReceiveMessage(msg slack.Message) bool {
	if p.config.CheckEnabledMessage(msg) {
		prefix := fmt.Sprintf("@%s %s", p.config.MentionName, p.config.GetSubcommand("calendar"))
		if strings.HasPrefix(msg.Text, prefix) {
			go func(m slack.Message) {
				p.sendMessage(msg.Channel)
			}(msg)
			return true
		}
	}
	return false
}

func (p *Plugin) sendMessage(channel string) {
	min, max := getTimeMinMax()
	items := p.calendar.GetCalendarItems(min, max, p.getSecretPath(), p.getTokenPath())
	fields := make([]slack.AttachmentField, 0)
	for _, i := range items {
		if i.isExcluded(p.config) {
			continue
		}
		start := fmt.Sprintf("%02d:%02d", i.Start.Hour(), i.Start.Minute())
		end := fmt.Sprintf("%02d:%02d", i.End.Hour(), i.End.Minute())
		field := createAttachmentField(i.Summary, fmt.Sprintf("%s〜%s / %s", start, end, i.Location))
		fields = append(fields, field)
	}
	if len(fields) > 0 {
		attachment := createAttachment("今日の予定です!", fields)
		p.client.SendAttachment(p.client.GetBotName(), attachment, channel)
	} else {
		p.client.SendMessage("今日の予定はないです...", channel)
	}
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

func createAttachment(pretext string, fields []slack.AttachmentField) slack.Attachment {
	return slack.Attachment{
		Pretext: pretext,
		Color:   "#3e6cf7",
		Fields:  fields,
	}
}

func createAttachmentField(title string, value string) slack.AttachmentField {
	return slack.AttachmentField{
		Title: title,
		Value: value,
	}
}

// Calendar for calendar repository interface
type Calendar interface {
	GetCalendarItems(min time.Time, max time.Time, secretPath string, tokenPath string) []Item
}

// Generator is calendar plugin generator
type Generator struct {
	Calendar Calendar
}

// Item is one of calendar event
type Item struct {
	Summary  string
	Location string
	Start    time.Time
	End      time.Time
}

// Generate generates Calendar plugin
func (g *Generator) Generate(config plugin.Config, client slack.Client) plugin.Plugin {
	plugin := &Plugin{
		config:   config,
		client:   client,
		calendar: g.Calendar,
	}
	scheduler := &schedule.Scheduler{Client: client, Config: config}
	scheduler.Start(func(c string) { plugin.sendMessage(c) })
	return plugin
}

func getTimeMinMax() (time.Time, time.Time) {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return now, tomorrow
}

func (item *Item) isExcluded(config plugin.Config) bool {
	exclude, ok := config.Data["exclude"]
	if !ok {
		return false
	}
	for _, e := range exclude.([]interface{}) {
		if item.Summary == e.(string) {
			return true
		}
	}
	return false
}
