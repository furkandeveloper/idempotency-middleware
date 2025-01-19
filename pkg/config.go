package pkg

import "time"

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
