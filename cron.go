package main

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

func scheduleIncreaseCap(ctx context.Context, g *golim) {
	userKeys := g.cache.getAllUserLimitersKeys(ctx)
	cr := cron.New()
	_, err := cr.AddFunc("@every 1m", func() {
		fmt.Println("Running tasks")
		for _, key := range userKeys {
			g.cache.increaseCap(ctx, key, g.limiterRole)
		}
	})
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}
	cr.Start()
}
