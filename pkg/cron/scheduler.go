package scheduler

import "github.com/robfig/cron/v3"

func Schedule(spec string, f func()) {
	c := cron.New()
	c.AddFunc(spec, f)
	c.Start()
}
