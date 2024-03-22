package main

import (
	"context"
	"encoding/json"
	"github.com/khalil-farashiani/golim/role"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type cache struct {
	*redis.Client
}

func initRedis() *cache {
	url := os.Getenv("REDIS_URI")
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return &cache{
		redis.NewClient(opts),
	}
}

func (c *cache) increaseCap(ctx context.Context, rl *limiterRole, amount int64) {
	key := operationIdToString[rl.operation] + " " + rl.endPoint
	c.Do(ctx, "INCRBY", key, amount)
}

func (c *cache) decreaseCap(ctx context.Context, rl *limiterRole) {
	key := operationIdToString[rl.operation] + " " + rl.endPoint
	c.Do(ctx, "DECR", key)
}

func (c *cache) setLimiter(ctx context.Context, params *role.GetRoleParams, val *role.GetRoleRow) {
	key := params.Operation + " " + params.Endpoint
	err := c.Set(ctx, key, val, time.Minute*60).Err()
	if err != nil {
		panic(err)
	}
}

func (c *cache) getLimiter(ctx context.Context, params role.GetRoleParams) *role.GetRoleRow {
	var res role.GetRoleRow
	var key = params.Operation + " " + params.Endpoint
	val, err := c.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	json.Unmarshal([]byte(val), &res)
	return &res
}

func (c *cache) getUserRequestCap(ctx context.Context, ipAddr string, rl *limiterRole) int64 {
	key := ipAddr + rl.endPoint + rl.endPoint
	var res = new(int64)
	val, err := c.Get(ctx, key).Result()

	if err != nil && err != redis.Nil {
		panic(err)
	}
	json.Unmarshal([]byte(val), res)
	return *res
}
