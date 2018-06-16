package plugin

import (
	"fmt"
	"log"

	"github.com/naokirin/slan-go/app/domain/response"
	"github.com/naokirin/slan-go/app/infrastructure/yaml"

	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/infrastructure/slack"
)

// GeneratePluginProcessArgs is arguments of GeneratePluginProcess function
type GeneratePluginProcessArgs struct {
	Client           *slack.Client
	MentionName      string
	Language         string
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
			MentionName:       args.MentionName,
			ResponseTemplates: &response.Template{},
			Data:              v.(map[interface{}]interface{}),
		}
		m, err := yaml.ParseFromFile(fmt.Sprintf("responses/%s/%s.yaml", args.Language, pluginName))
		if err != nil {
			log.Panicf("error: %v", err)
		}
		responses := make(map[string]string)
		responseTemplate, ok := pluginConfig.Data["response_template"]
		if ok {
			m, err = yaml.ParseFromFile(fmt.Sprintf("responses/%s.yaml", responseTemplate.(string)))
			if err != nil {
				log.Panicf("error: %v", err)
			}
		}
		for k, v := range m {
			responses[k.(string)] = v.(string)
		}
		pluginConfig.ResponseTemplates.AddTemplates(responses)
		result = append(result, pg.Generate(pluginConfig, args.Client))
	}
	return result
}
