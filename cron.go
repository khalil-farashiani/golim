package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

var cr = cron.New()

func runCronTasks(ctx context.Context, g *golim) {
	_, err := cr.AddFunc("@every 1m", func() {
		scheduleIncreaseCap(ctx, g)
	})
	if err != nil {
		g.logger.errLog.Println(err)
		fmt.Println("Error scheduling task:", err)
	}
	cr.Start()
}

func scheduleIncreaseCap(ctx context.Context, g *golim) {
	roles, err := g.getRoles(ctx)
	if err != nil {
		g.logger.errLog.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, role := range roles {
		userKeys := g.cache.getAllUserLimitersKeys(ctx, limiterCacheRegexPatternKey+role.Operation+role.Endpoint)
		for _, key := range userKeys {
			wg.Add(1)
			go func(ctx context.Context, key string, tokenAmount int64) {
				defer wg.Done()
				g.cache.increaseCap(ctx, key, tokenAmount)
			}(context.Background(), key, g.limiterRole.addToken)
		}
	}
	wg.Wait()
}
