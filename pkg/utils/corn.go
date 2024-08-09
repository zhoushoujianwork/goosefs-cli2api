package utils

import (
	"github.com/robfig/cron"
)

func StartScheduler() {
	c := cron.New()

	// 每 1 分钟执行一次
	c.AddFunc("0 * * * * *", func() {
	})

	c.Start()
	defer c.Stop()
}
