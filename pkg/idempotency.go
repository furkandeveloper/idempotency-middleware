package pkg

import (
	"bytes"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Config struct {
	Redis       RedisConfig
	Idempotency IdempotencyConfig
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type IdempotencyConfig struct {
	HeaderKey      string
	ExpirationTime time.Duration
}

func LoadConfig(expirationTime int, headerKey string) Config {
	return Config{
		Redis: RedisConfig{
			Address:  "localhost:6379",
			Password: "",
			DB:       0,
		},
		Idempotency: IdempotencyConfig{
			HeaderKey:      headerKey,
			ExpirationTime: time.Duration(expirationTime) * time.Minute,
		},
	}
}

type Option struct {
	HeaderKey      string
	ExpirationTime time.Duration
}

type responseRecorder struct {
	writer http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *responseRecorder) Header() http.Header {
	return r.writer.Header()
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.writer.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b) // Yanıtı kaydet
	return r.writer.Write(b)
}

func NewIdempotencyMiddleware(redisClient *redis.Client, opt Option) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(opt.HeaderKey)
			if requestID == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Missing idempotency key")
			}

			ctx := context.Background()
			if cachedResponse, err := redisClient.Get(ctx, requestID).Result(); err == nil {
				return c.String(http.StatusOK, cachedResponse)
			}

			rec := &responseRecorder{
				writer: c.Response().Writer,
			}
			c.Response().Writer = rec

			if err := next(c); err != nil {
				c.Error(err)
				return err
			}

			redisClient.Set(ctx, requestID, rec.body.String(), opt.ExpirationTime)

			return nil
		}
	}
}
