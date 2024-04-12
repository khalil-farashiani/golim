package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/khalil-farashiani/golim/internal/domain"
	"github.com/khalil-farashiani/golim/internal/store/role"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	limiterCacheMainKey         = "GOLIM_KEY"
	limiterCacheRegexPatternKey = "*GOLIM_KEY"
)

const (
	initialTokenForTheFirstUSerRequest = 1
)

// Cache struct to handle redis operations
type Cache struct {
	*redis.Client
}

// InitRedis initializes redis connection
func InitRedis() *Cache {
	url := os.Getenv("REDIS_URI")
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Error parsing Redis URL: %v", err)
	}
	return &Cache{
		redis.NewClient(opts),
	}
}

// getAllUserLimitersKeys retrieves all user limiter keys from Cache
func (c *Cache) getAllUserLimitersKeys(ctx context.Context, pattern string) []string {
	res, err := c.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("Error retrieving keys: %v", err)
		return nil
	}
	return res
}

// increaseCap increases the capacity in Cache for a given key
func (c *Cache) increaseCap(ctx context.Context, key string, tokenAmount int64) {
	if err := c.IncrBy(ctx, key, tokenAmount).Err(); err != nil {
		log.Printf("Error increasing capacity: %v", err)
	}
}

// decreaseCap decreases the capacity in Cache for a given key
func (c *Cache) decreaseCap(ctx context.Context, userIP string, rl *domain.LimiterRole) {
	key := fmt.Sprintf("%s%s%s%s", userIP, limiterCacheMainKey, rl.Operation, rl.EndPoint)
	if err := c.Decr(ctx, key).Err(); err != nil {
		log.Printf("Error decreasing capacity: %v", err)
	}
}

// setLimiter sets a limiter in Cache based on parameters
func (c *Cache) setLimiter(ctx context.Context, params *role.GetRoleParams, val *role.GetRoleRow) {
	key := fmt.Sprintf("%s %s", params.Operation, params.Endpoint)
	if err := c.Set(ctx, key, *val, time.Minute*60).Err(); err != nil {
		log.Printf("Error setting limiter: %v", err)
	}
}

// getLimiter retrieves a limiter from Cache based on parameters
func (c *Cache) getLimiter(ctx context.Context, params role.GetRoleParams) *role.GetRoleRow {
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

// setUserRequestCap sets user request capacity in Cache
func (c *Cache) setUserRequestCap(ctx context.Context, key string, role role.GetRoleRow) {
	if err := c.Set(ctx, key, role.InitialTokens, time.Hour).Err(); err != nil {
		log.Printf("Error setting user request capacity: %v", err)
	}
}

// getUserRequestCap retrieves user request capacity from Cache
func (c *Cache) getUserRequestCap(ctx context.Context, ipAddr string, g *domain.Golim, role role.GetRoleRow) int64 {
	key := fmt.Sprintf("%s%s%s %s", ipAddr, limiterCacheMainKey, g.LimiterRole.Operation, g.LimiterRole.EndPoint)
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
