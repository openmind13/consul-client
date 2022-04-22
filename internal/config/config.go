package config

import (
	"consul-client/internal/consul"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var validate = validator.New()

type Config struct {
	Consul consul.Config `mapstructure:"consul"`
}

func (c Config) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}
	return nil
}

func Get(path string) (Config, error) {
	c := Config{}
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}
	return c, nil
}
