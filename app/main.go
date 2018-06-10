package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/naokirin/slan-go/app/application/plugin"
	"github.com/naokirin/slan-go/app/domain/calendar"
	"github.com/naokirin/slan-go/app/domain/lgtmize"
	"github.com/naokirin/slan-go/app/domain/lunch"
	"github.com/naokirin/slan-go/app/domain/memolist"
	"github.com/naokirin/slan-go/app/domain/ping"
	dplugin "github.com/naokirin/slan-go/app/domain/plugin"
	icalendar "github.com/naokirin/slan-go/app/infrastructure/google/calendar"
	"github.com/naokirin/slan-go/app/infrastructure/google/spreadsheets"
	ilgtmize "github.com/naokirin/slan-go/app/infrastructure/lgtmize"
	"github.com/naokirin/slan-go/app/infrastructure/slack"
	imemolist "github.com/naokirin/slan-go/app/infrastructure/sqlite/memolist"
	"github.com/naokirin/slan-go/app/infrastructure/yaml"
)

var pluginGenerators = map[string]dplugin.Generator{
	"memolist": &memolist.Generator{Repository: &imemolist.Memo{}},
	"ping":     &ping.Generator{},
	"calendar": &calendar.Generator{Calendar: &icalendar.Calendar{}},
	"lunch":    &lunch.Generator{Repository: &spreadsheets.Spreadsheets{}},
	"lgtmize":  &lgtmize.Generator{LGTMize: &ilgtmize.LGTMize{}},
}

func main() {
	log.Println("Start slan-go")
	config := yaml.GetConfigurationRepository()
	location, err := time.LoadLocation(config.GetLocation())
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		time.Local = location
	}

	pluginConfigs := config.GetPlugins()
	client := slack.CreateClient(config.GetSlackToken())
	plugins := plugin.GeneratePlugins(plugin.GeneratePluginProcessArgs{
		Client:           client,
		MentionName:      config.GetMentionName(),
		PluginConfigs:    pluginConfigs,
		PluginGenerators: pluginGenerators,
	})
	rand.Seed(time.Now().UnixNano())
	defaultResponses := config.GetDefaultResponses()
	defaultResponsesLen := len(defaultResponses)
	for msg := range client.GenerateReceivedEventChannel() {
		match := false
		for _, p := range plugins {
			match = p.ReceiveMessage(msg) || match
		}
		if !match && strings.HasPrefix(msg.Text, "@"+config.GetMentionName()) && defaultResponsesLen > 0 {
			n := rand.Intn(defaultResponsesLen)
			client.SendMessage(defaultResponses[n], msg.Channel)
		}
	}
}
