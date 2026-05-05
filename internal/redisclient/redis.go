package redisclient

// import (
// 	"context"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// type Client struct {
// 	RDB *redis.Client
// }

// func NewRedisClient(addr string) *Client {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:         addr,
// 		DB:           0,
// 		DialTimeout:  5 * time.Minute,
// 		ReadTimeout:  3 * time.Minute,
// 		WriteTimeout: 3 * time.Minute,
// 	})

// 	return &Client{RDB: rdb}
// }

// func (c *Client) Ping(ctx context.Context) error {
// 	return c.RDB.Ping(ctx).Err()
// }
