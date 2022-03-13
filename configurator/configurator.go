package configurator

import (
	"consul-client/config"

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
		newConf := config.Config{}
		if err := viper.Unmarshal(&newConf); err != nil {
			ConfigErrorChan <- err
		} else {
			if err := newConf.Validate(); err != nil {
				ConfigErrorChan <- err
			} else {
				ConfigChan <- newConf
			}
		}
	})

	return nil
}

// func StartConsulWatcher() {
// 	consulWatcher := viper.New()
// 	// consulWatcher.AddRemot
// 	shedu
// }
