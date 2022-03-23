package config

import "fmt"

type Config struct {
	CommonConfig           CommonConfig           `mapstructure:"common"`
	ServiceDiscoveryConfig ServiceDiscoveryConfig `mapstructure:"service_discovery"`
}

func (c *Config) Validate() error {
	if err := c.CommonConfig.Validate(); err != nil {
		return err
	}
	if err := c.ServiceDiscoveryConfig.Validate(); err != nil {
		return err
	}
	return nil
}

type CommonConfig struct {
	LogLevel string `mapstructure:"log_level"`
}

var (
	commonConfigErrorTemplate = "Field '%s' in common config is empty"
)

func (cc *CommonConfig) Validate() error {
	if cc.LogLevel == "" {
		return fmt.Errorf(commonConfigErrorTemplate, "log_level")
	}
	return nil
}

type ServiceDiscoveryConfig struct {
	Addr           string `mapstructure:"addr"`
	Token          string `mapstructure:"token"`
	HttpListenAddr string `mapstructure:"http_listen_addr"`
}

var (
	serviceDiscoveryConfigErrorTemplate = "Field '%s' in service discovery config is empty"
)

func (sdcfg *ServiceDiscoveryConfig) Validate() error {
	if sdcfg.Addr == "" {
		return fmt.Errorf(serviceDiscoveryConfigErrorTemplate, "addr")
	}
	if sdcfg.Token == "" {
		return fmt.Errorf(serviceDiscoveryConfigErrorTemplate, "token")
	}
	if sdcfg.HttpListenAddr == "" {
		return fmt.Errorf(serviceDiscoveryConfigErrorTemplate, "http_addr")
	}
	return nil
}
