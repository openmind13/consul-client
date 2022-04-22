package main

import (
	"consul-client/internal/config"
	"consul-client/internal/consul"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	cfgPath = flag.String("cfg_path", "./config.toml", "config path")
)

const (
	version = "v0.0.0-dev"
)

func main() {
	flag.Parse()

	errChan := make(chan error, 1)

	conf, err := config.Get(*cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := conf.Validate(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Creating consul service")
	consul, err := consul.NewClient(conf.Consul, version)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Starting consul service")
	go consul.Start(errChan)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		consul.StopService()
		logrus.Fatal(err)
	case sig := <-signalChan:
		logrus.Info("Killed by signal: ", sig)
		consul.StopService()
		os.Exit(0)
	}
}
