package plugin

import (
	"log"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/infrastructure/slack"
)

// GeneratePluginProcessArgs is arguments of GeneratePluginProcess function
type GeneratePluginProcessArgs struct {
	Client           *slack.Client
	MentionName      string
	PluginConfigs    []interface{}
	PluginGenerators map[string]plugin.Generator
}

// GeneratePlugins runs plugin goroutines
func GeneratePlugins(args GeneratePluginProcessArgs) []plugin.Plugin {
	result := make([]plugin.Plugin, 0)
	for _, v := range args.PluginConfigs {
		pName, ok := v.(map[interface{}]interface{})["plugin"]
		if !ok {
			continue
		}
		pluginName := pName.(string)
		pg, ok := args.PluginGenerators[pluginName]
		if !ok {
			log.Printf("Plugin: %s is not found\n", pluginName)
			continue
		}
		pluginConfig := plugin.Config{
			MentionName: args.MentionName,
			Data:        v.(map[interface{}]interface{}),
		}
		result = append(result, pg.Generate(pluginConfig, args.Client))
	}
	return result
}
