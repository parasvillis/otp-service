package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
	Exists(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HSet(ctx context.Context, key string, values map[string]interface{}) *redis.IntCmd
	Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error)
	RunLuaAny(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	RunLuaInt(ctx context.Context, script string, keys []string, args ...interface{}) (int, error)
}
type cache struct {
	client *redis.Client
}

func NewCache(c *redis.Client) Cache {
	return &cache{client: c}
}
func (r *cache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *cache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}

func (r *cache) Exists(ctx context.Context, key string) (int64, error) {
	return r.client.Exists(ctx, key).Result()
}

func (r *cache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (c *cache) HSet(ctx context.Context, key string, values map[string]interface{}) *redis.IntCmd {
	return c.client.HSet(ctx, key, values)
}

func (c *cache) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return c.client.HGet(ctx, key, field)
}

func (c *cache) Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd {
	return c.client.Expire(ctx, key, ttl)
}

func (r *cache) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return r.client.HIncrBy(ctx, key, field, incr).Result()
}

func (r *cache) RunLuaAny(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	luaScript := redis.NewScript(script)
	cmd := luaScript.Run(ctx, r.client, keys, args...)
	return cmd.Result()
}

func (r *cache) RunLuaInt(ctx context.Context, script string, keys []string, args ...interface{}) (int, error) {
	res, err := r.RunLuaAny(ctx, script, keys, args...)
	if err != nil {
		return 0, err
	}
	i, ok := res.(int64)
	if !ok {
		return 0, fmt.Errorf("expected int64 from Lua, got %T", res)
	}
	return int(i), nil
}
