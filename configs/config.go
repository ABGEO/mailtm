package configs

import "time"

type Config struct {
	APIBaseURL       string
	APIClientTimeout time.Duration
}

func NewConfig() *Config {
	return &Config{
		APIBaseURL:       "https://api.mail.tm",
		APIClientTimeout: 30 * time.Second,
	}
}
