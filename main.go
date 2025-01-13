package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"idempotency-middleware/internal/handlers"
	idempotency "idempotency-middleware/internal/middleware"
	"idempotency-middleware/internal/redis"
)

// main function
func main() {
	cfg := redis.LoadConfig(1, "X-Request-Id")

	redisClient := redis.NewClient(cfg.Redis)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	idempotencyMiddleware := idempotency.NewIdempotencyMiddleware(redisClient, idempotency.Config(cfg.Idempotency))

	e.GET("/no-middleware", handlers.ExampleHandlerWithoutIdempotency)
	e.GET("/with-middleware", handlers.ExampleHandler, idempotencyMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
