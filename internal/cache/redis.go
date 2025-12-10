package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

var Rdc *redis.Client

func InitRedisConnection() error {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		log.Fatal("Addr to redis is not passed!")
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		log.Fatal("Addr to redis is not passed!")
	}
	Rdc = redis.NewClient(&redis.Options{
		Addr:         host + port,
		Password:     "",
		DB:           0,
		PoolSize:     20,
		MinIdleConns: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return Rdc.Ping(ctx).Err()
}
