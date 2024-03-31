package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/khalil-farashiani/golim/role"
	"github.com/redis/go-redis/v9"
)

// cache struct to handle redis operations
type cache struct {
	*redis.Client
}

// initRedis initializes redis connection
func initRedis() *cache {
	url := os.Getenv("REDIS_URI")
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Error parsing Redis URL: %v", err)
	}
	return &cache{
		redis.NewClient(opts),
	}
}

// getAllUserLimitersKeys retrieves all user limiter keys from cache
func (c *cache) getAllUserLimitersKeys(ctx context.Context) []string {
	res, err := c.Keys(ctx, limiterCacheRegexPatternKey).Result()
	if err != nil {
		log.Printf("Error retrieving keys: %v", err)
		return nil
	}
	return res
}

// increaseCap increases the capacity in cache for a given key
func (c *cache) increaseCap(ctx context.Context, key string, rl *limiterRole) {
	if err := c.IncrBy(ctx, key, rl.addToken).Err(); err != nil {
		log.Printf("Error increasing capacity: %v", err)
	}
}

// decreaseCap decreases the capacity in cache for a given key
func (c *cache) decreaseCap(ctx context.Context, userIP string, rl *limiterRole) {
	key := fmt.Sprintf("%s%s%s %s", userIP, limiterCacheMainKey, rl.operation, rl.endPoint)
	if err := c.Decr(ctx, key).Err(); err != nil {
		log.Printf("Error decreasing capacity: %v", err)
	}
}

// setLimiter sets a limiter in cache based on parameters
func (c *cache) setLimiter(ctx context.Context, params *role.GetRoleParams, val *role.GetRoleRow) {
	key := fmt.Sprintf("%s %s", params.Operation, params.Endpoint)
	if err := c.Set(ctx, key, *val, time.Minute*60).Err(); err != nil {
		log.Printf("Error setting limiter: %v", err)
	}
}

// getLimiter retrieves a limiter from cache based on parameters
func (c *cache) getLimiter(ctx context.Context, params role.GetRoleParams) *role.GetRoleRow {
	var res role.GetRoleRow
	key := fmt.Sprintf("%s %s", params.Operation, params.Endpoint)
	val, err := c.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("Error getting limiter: %v", err)
		}
		return nil
	}
	json.Unmarshal([]byte(val), &res)
	return &res
}

// setUserRequestCap sets user request capacity in cache
func (c *cache) setUserRequestCap(ctx context.Context, key string, role role.GetRoleRow) {
	if err := c.Set(ctx, key, role.InitialTokens, time.Hour).Err(); err != nil {
		log.Printf("Error setting user request capacity: %v", err)
	}
}

// getUserRequestCap retrieves user request capacity from cache
func (c *cache) getUserRequestCap(ctx context.Context, ipAddr string, g *golim, role role.GetRoleRow) int64 {
	key := fmt.Sprintf("%s%s%s %s", ipAddr, limiterCacheMainKey, g.limiterRole.operation, g.limiterRole.endPoint)
	val, err := c.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("Error getting user request capacity: %v", err)
			return 0
		}
		go c.setUserRequestCap(ctx, key, role)
		return initialTokenForTheFirstUSerRequest
	}
	var res int64
	json.Unmarshal([]byte(val), &res)
	return res
}
