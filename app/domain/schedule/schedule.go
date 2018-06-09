package schedule

import (
	"github.com/naokirin/slan-go/app/domain/plugin"
	"github.com/naokirin/slan-go/app/domain/slack"
	"github.com/robfig/cron"
)

// Scheduler for scheduled plugin
type Scheduler struct {
	Client slack.Client
	Config plugin.Config
}

// Start is starting scheduler
func (scheduler *Scheduler) Start(callback func(channel string)) {
	config, ok := scheduler.Config.Data["schedule"]
	if !ok {
		return
	}
	sc, ok := config.(map[interface{}]interface{})
	if !ok {
		return
	}
	ch, ok := sc["channel"]
	if !ok {
		return
	}
	ss, ok := sc["cron"]
	if !ok {
		return
	}
	channel, ok := scheduler.Client.ConvertChannelNameToID(ch.(string))
	if !ok {
		return
	}
	c := cron.New()
	s := ss.(map[interface{}]interface{})
	e, ok := s["expr"]
	if !ok {
		return
	}
	expr := e.(string)
	f := func() { callback(channel) }
	sk, sok := s["skip"]
	o, ook := s["offset"]
	if sok || ook {
		skip := 0
		offset := 0
		if sok {
			skip = sk.(int)
		}
		if ook {
			offset = o.(int)
		}
		i := 0
		f = func() {
			if i >= offset && (i-offset)%skip == 0 {
				callback(channel)
			}
			i++
		}
	}
	c.AddFunc(expr, f)
	c.Start()
}
