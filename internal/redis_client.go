package internal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(redisAddr string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr, // Redis sunucusunun adresi (localhost:6379 vb.)
	})
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis bağlantısı kurulamadı: %v", err)
	}
	fmt.Println("Redis bağlantısı başarılı")
}
