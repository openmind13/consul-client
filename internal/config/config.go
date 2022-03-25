package config

type Config struct {
	ConsulAddress     string `mapstructure:"consul_address" validate:"required"`
	Token             string `mapstructure:"token" validate:"required"`
	HttpListenAddress string `mapstructure:"http_listen_address"`
}
