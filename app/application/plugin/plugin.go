package plugin

import (
	"fmt"

	"github.com/naokirin/slan-go/app/domain/memolist"
	"github.com/naokirin/slan-go/app/domain/ping"
	"github.com/naokirin/slan-go/app/domain/plugin"
	dslack "github.com/naokirin/slan-go/app/domain/slack"
	"github.com/naokirin/slan-go/app/infrastructure/slack"
	imemolist "github.com/naokirin/slan-go/app/infrastructure/sqlite/memolist"
)

var plugins = map[string]func(plugin.Config, *slack.Client, chan dslack.Message){
	"memolist": func(config plugin.Config, client *slack.Client, in chan dslack.Message) {
		memolist.GeneratePluginGoroutine(config, imemolist.Memo{}, generateSender(client), in)
	},
	"ping": func(config plugin.Config, client *slack.Client, in chan dslack.Message) {
		ping.GeneratePluginGoroutine(config, generateSender(client), in)
	},
}

func generateSender(client *slack.Client) func(string, string) {
	return func(text string, channel string) { client.SendMessage(text, channel) }
}

// GeneratePluginProcessArgs is arguments of GeneratePluginProcess function
type GeneratePluginProcessArgs struct {
	Client        *slack.Client
	MentionName   string
	PluginConfigs []interface{}
}

// GeneratePluginProcess runs plugin goroutines
func GeneratePluginProcess(args GeneratePluginProcessArgs) []chan dslack.Message {
	chans := make([]chan dslack.Message, 2)
	for _, v := range args.PluginConfigs {
		out := make(chan dslack.Message)
		chans = append(chans, out)
		pName, ok := v.(map[interface{}]interface{})["plugin"]
		if !ok {
			return chans
		}
		pluginName := pName.(string)
		p, ok := getPluginGenerator(pluginName)
		if !ok {
			fmt.Printf("Plugin: %s is not found", pluginName)
			return chans
		}
		pluginConfig := plugin.Config{
			MentionName: args.MentionName,
			Data:        v.(map[interface{}]interface{}),
		}
		p(pluginConfig, args.Client, out)
	}
	return chans
}

func getPluginGenerator(pluginName string) (func(plugin.Config, *slack.Client, chan dslack.Message), bool) {
	v, ok := plugins[pluginName]
	return v, ok
}
