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
	if err := configurator.Validator.Struct(config); err != nil {
		log.Fatal(err)
	}

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

	serviceDiscovery, err := servicediscovery.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	serviceDiscovery.DeregisterService()

	if err := serviceDiscovery.RegisterService(); err != nil {
		log.Fatal(err)
	}
	log.Println("Consul service", servicediscovery.SERVICE_NAME, "registered")

	go serviceDiscovery.StartServiceUpdater()

	// go serviceDiscovery.StartService(errChan)
	// go serviceDiscovery.ServiceListenHttp(errChan)

	serviceDiscovery.DeregisterService()

	err = <-errChan
	serviceDiscovery.DeregisterService()
	log.Println(err)
}
