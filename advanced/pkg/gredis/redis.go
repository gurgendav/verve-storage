package gredis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var RedisConn *redis.Client

var ctx = context.Background()

func init() {
	redisAddr := os.Getenv("REDIS_ADDR")

	RedisConn = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

}

func Set(key string, data interface{}, expiration time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return RedisConn.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	reply, err := RedisConn.Get(ctx, key).Result()
	return reply, err
}

func Delete(key string) (bool, error) {
	suc, err := RedisConn.Del(ctx, key).Result()

	return suc > 0, err
}

func DropDatabase() error {
	return RedisConn.FlushDB(ctx).Err()
}
