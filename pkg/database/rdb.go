package database

import (
	"fmt"
	"github.com/go-redis/redis/v9"
)

type RedisConnectionInfo struct {
	Host     string
	Port     int
	Password string
}

func NewRedisClient(rdbInfo RedisConnectionInfo) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdbInfo.Host, rdbInfo.Port),
		Password: rdbInfo.Password,
		DB:       0,
	})
}
