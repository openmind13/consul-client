package main

import (
	"consul-client/configurator"
	"consul-client/servicediscovery"
	"flag"
	"log"
)

var (
	cfgPath = flag.String("config", "./config", "config path")
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
				log.Println("Error in parsing toml config: ", configErr)
			}
		}
	}()

	serviceDiscovery, err := servicediscovery.NewClient(config.ServiceDiscoveryConfig)
	if err != nil {
		log.Fatal(err)
	}

	go serviceDiscovery.Start(errChan)
	go serviceDiscovery.Listen(errChan)

	err = <-errChan
	log.Println(err)
}
