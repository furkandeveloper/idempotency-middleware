package middleware

import (
	"fmt"
	"github.com/furkandeveloper/idempotency-middleware/internal"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func IdempotencyMiddleware(redisClient *redis.Client, config *internal.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// X-Request-Id header'ını al
			requestID := c.Request().Header.Get("X-Request-Id")
			if requestID == "" {
				// Eğer başlık yoksa, yeni bir request işlemine devam et
				return next(c)
			}

			// Redis cache'inde var mı kontrol et
			cacheKey := fmt.Sprintf("idempotency:%s", requestID)
			cachedResponse, err := redisClient.Get(c.Request().Context(), cacheKey).Result()
			if err == redis.Nil {
				// Eğer cache'de yoksa, normal işlemi yap ve yanıtı kaydet
				err := next(c)
				if err != nil {
					return err
				}

				// Yanıtı cache'le
				// Cache'e response'ı yerleştir
				// Burada örnek olarak yanıtı string olarak kaydedeceğiz. Gerçek hayatta JSON yanıt kaydedilebilir
				redisClient.Set(c.Request().Context(), cacheKey, "processed", config.ExpireDuration)
				return nil
			} else if err != nil {
				log.Printf("Redis hata: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
			}

			// Eğer cache'de varsa, mevcut response'u döndür
			c.String(http.StatusOK, cachedResponse)
			return nil
		}
	}
}
