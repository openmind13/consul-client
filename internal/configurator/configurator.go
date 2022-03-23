package configurator

import (
	"consul-client/internal/config"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

var (
	ConfigChan      = make(chan config.Config, 1)
	ConfigErrorChan = make(chan error, 1)
)

type Options struct{}

func StartTomlWatcher(configPath string) error {
	localWatcher := viper.New()
	localWatcher.SetConfigName("config.toml")
	localWatcher.SetConfigType("toml")
	localWatcher.AddConfigPath(".")

	if err := localWatcher.ReadInConfig(); err != nil {
		return err
	}

	var conf config.Config
	if err := localWatcher.Unmarshal(&conf); err != nil {
		return err
	}
	ConfigChan <- conf

	localWatcher.WatchConfig()

	localWatcher.OnConfigChange(func(in fsnotify.Event) {
		if err := localWatcher.ReadInConfig(); err != nil {
			ConfigErrorChan <- err
			return
		}
		newConf := config.Config{}
		if err := viper.Unmarshal(&newConf); err != nil {
			ConfigErrorChan <- err
		} else {
			fmt.Printf("%+v\n", newConf)
			if err := newConf.Validate(); err != nil {
				ConfigErrorChan <- err
			} else {
				ConfigChan <- newConf
			}
		}
	})

	return nil
}
