package main

import (
	"fmt"
	"github.com/furkandeveloper/idempotency-middleware/internal"
	"github.com/furkandeveloper/idempotency-middleware/middleware"
	"github.com/labstack/echo/v4"
	"os"
	"time"
)

func main() {
	// Redis bağlantısını başlat
	redisAddr := "localhost:6379" // Redis adresi
	internal.InitRedis(redisAddr)

	// Config ayarları (expire süresi 10 dakika)
	config := internal.NewConfig(10 * time.Minute)

	// Echo uygulamasını başlat
	e := echo.New()

	// Idempotency middleware'i uygulamaya dahil et
	e.Use(middleware.IdempotencyMiddleware(internal.RedisClient, config))

	// Basit bir test route'u
	e.GET("/example", func(c echo.Context) error {
		return c.String(200, "Request processed successfully")
	})

	// Uygulamayı başlat
	port := ":8080"
	fmt.Printf("Uygulama %s portunda başlatılıyor...\n", port)
	if err := e.Start(port); err != nil {
		fmt.Printf("Uygulama başlatılamadı: %v\n", err)
		os.Exit(1)
	}
}
