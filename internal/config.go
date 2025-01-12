package internal

import "time"

// Config, idempotency middleware yapılandırmasını içerir.
type Config struct {
	ExpireDuration time.Duration
}

func NewConfig(expireDuration time.Duration) *Config {
	return &Config{ExpireDuration: expireDuration}
}
