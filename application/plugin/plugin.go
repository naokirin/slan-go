package plugin

import (
	"fmt"

	"github.com/naokirin/slan-go/application/plugin/memolist"
	"github.com/naokirin/slan-go/domain/plugin"
	dslack "github.com/naokirin/slan-go/domain/slack"
	"github.com/naokirin/slan-go/infrastructure/slack"
)

var plugins = map[string]func(plugin.Config, *slack.Client, chan dslack.Message){
	"memolist": func(config plugin.Config, client *slack.Client, in chan dslack.Message) {
		memolist.GeneratePluginGoroutine(config, client, in)
	},
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