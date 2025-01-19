package main

import (
	"github.com/furkandeveloper/idempotency-middleware/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// main function
func main() {
	cfg := pkg.LoadConfig(1, "X-Request-Id")

	redisClient := pkg.NewClient(cfg.Redis)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	idempotencyMiddleware := pkg.NewIdempotencyMiddleware(redisClient, pkg.Option(cfg.Idempotency))

	e.GET("/no-middleware", pkg.ExampleHandlerWithoutIdempotency)
	e.GET("/with-middleware", pkg.ExampleHandler, idempotencyMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
