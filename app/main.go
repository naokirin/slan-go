package main

import (
	"github.com/naokirin/slan-go/app/application/plugin"
	"github.com/naokirin/slan-go/app/domain/memolist"
	"github.com/naokirin/slan-go/app/domain/ping"
	dplugin "github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/infrastructure/slack"
	imemolist "github.com/naokirin/slan-go/app/infrastructure/sqlite/memolist"
	"github.com/naokirin/slan-go/app/infrastructure/yaml"
)

var pluginGenerators = map[string]dplugin.Generator{
	"memolist": &memolist.Generator{Repository: imemolist.Memo{}},
	"ping":     &ping.Generator{},
}

func main() {
	config := yaml.GetConfigurationRepository()
	pluginConfigs := config.GetPlugins()
	client := slack.CreateClient(config.GetSlackToken())
	plugins := plugin.GeneratePlugins(plugin.GeneratePluginProcessArgs{
		Client:           client,
		MentionName:      config.GetMentionName(),
		PluginConfigs:    pluginConfigs,
		PluginGenerators: pluginGenerators,
	})
	for msg := range client.GenerateReceivedEventChannel() {
		for _, p := range plugins {
			p.ReceiveMessage(msg)
		}
	}
}
