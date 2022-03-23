package main

import (
	"consul-client/internal/configurator"
	"consul-client/internal/servicediscovery"
	"flag"
	"log"
)

var (
	cfgPath = flag.String("cfg_path", "./config.toml", "config path")
)

func main() {
	flag.Parse()

	errChan := make(chan error, 1)

	if err := configurator.StartTomlWatcher(*cfgPath); err != nil {
		log.Fatal(err)
	}
	config := <-configurator.ConfigChan

	go func() {
		for {
			select {
			case config := <-configurator.ConfigChan:
				log.Printf("new config: %+v\n", config)
			case configErr := <-configurator.ConfigErrorChan:
				log.Println("Error in parsing toml config:", configErr)
			}
		}
	}()

	serviceDiscovery, err := servicediscovery.NewClient(config.ServiceDiscoveryConfig)
	if err != nil {
		log.Fatal(err)
	}

	go serviceDiscovery.Start(errChan)
	go serviceDiscovery.ServiceListenHttp(errChan)

	err = <-errChan
	log.Println(err)
}
