package configs

type Config struct {
	APIBaseURL string
}

func NewConfig() *Config {
	return &Config{
		APIBaseURL: "https://api.mail.tm",
	}
}
