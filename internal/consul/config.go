package consul

import (
	"time"
)

type Config struct {
	Address             string        `mapstructure:"address" validate:"required"`
	Token               string        `mapstructure:"token" validate:"required"`
	ServiceName         string        `mapstructure:"service_name" validate:"required"`
	ServiceID           string        `mapstructure:"service_id" validate:"required"`
	CheckName           string        `mapstructure:"check_name" validate:"required"`
	CheckTTL            time.Duration `mapstructure:"check_ttl" validate:"required"`
	UpdateCheckPeriod   time.Duration `mapstructure:"update_check_period" validate:"required"`
	ServiceUpdatePeriod time.Duration `mapstructure:"service_update_period" validate:"required"`
	Tags                []string      `mapstructure:"tags" validate:"required"`
}

func (c *Client) StartConfigWatcher() {
	// kvPair := consulapi.KVPair{
	// 	Key: "calc-config",
	// }
	// c.kvStorage.Acquire()

	// lockOpts := &consulapi.LockOptions{
	// 	Key: "calc-config",
	// }

	// fmt.Println("start config watcher")
	// for {
	// 	kvPair, _, err := c.kvStorage.Get(calcConfigName, &consulapi.QueryOptions{})
	// 	if err != nil {
	// 		logrus.Error(err)
	// 	} else {
	// 		// fmt.Println(kvPair.Key, string(kvPair.Value))

	// 		buf := bytes.Buffer{}
	// 		buf.Write(kvPair.Value)

	// 		// var config config.Config
	// 		// if err := json.NewDecoder(&buf).Decode(&config); err != nil {
	// 		// 	logrus.Error(err)
	// 		// }

	// 		// fmt.Printf("%+v\n", config)
	// 	}

	// 	time.Sleep(time.Second)
	// }

	select {}
}
