package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/khalil-farashiani/golim/role"
	"github.com/pingcap/log"
	"github.com/redis/go-redis/v9"
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

func (c *cache) getAllUserLimitersKeys(ctx context.Context) []string {
	res, err := c.Keys(ctx, "*GLOLIM_KEY*").Result()
	if err != nil {
		log.Fatal(err.Error())
	}
	return res
}

func (c *cache) increaseCap(ctx context.Context, key string, rl *limiterRole) {
	c.IncrBy(ctx, key, rl.addToken)
}

func (c *cache) decreaseCap(ctx context.Context, userIP string, rl *limiterRole) {
	key := userIP + "GLOLIM_KEY" + rl.operation + " " + rl.endPoint
	c.Decr(ctx, key)
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
