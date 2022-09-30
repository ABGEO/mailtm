package configs

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"github.com/spf13/viper"
)

type AuthConfig struct {
	ID    string
	Email string
	Token string
}

type Config struct {
	Auth struct {
		AuthConfig `mapstructure:",squash"`
	}
}

func NewConfig() (conf Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.mailtm")

	if err := viper.ReadInConfig(); err != nil {
		if ok := errors.As(err, &viper.ConfigFileNotFoundError{}); ok {
			if homedir, err := os.UserHomeDir(); err == nil {
				_ = os.Mkdir(homedir+"/.mailtm", os.ModePerm)
				_, _ = os.Create(homedir + "/.mailtm/config")
				_ = viper.WriteConfig()
			}
		}
	}

	_ = viper.Unmarshal(&conf)

	return conf
}

func (conf *Config) Write() {
	if cfg, err := json.Marshal(conf); err == nil {
		_ = viper.ReadConfig(bytes.NewBuffer(cfg))
		_ = viper.WriteConfig()
	}
}
