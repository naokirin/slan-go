package main

import (
	"github.com/naokirin/slan-go/app/application/plugin"
	"github.com/naokirin/slan-go/app/infrastructure/slack"
	"github.com/naokirin/slan-go/app/infrastructure/yaml"
)

func main() {
	config := yaml.GetConfigurationRepository()
	pluginConfigs := config.GetPlugins()
	client := slack.CreateClient(config.GetSlackToken())
	chans := plugin.GeneratePluginProcess(plugin.GeneratePluginProcessArgs{
		Client:        client,
		MentionName:   config.GetMentionName(),
		PluginConfigs: pluginConfigs,
	})
	for msg := range client.GenerateReceivedEventChannel() {
		for i := 0; i < len(chans); i++ {
			go func(index int) {
				chans[index] <- msg
			}(i)
		}
	}
}
