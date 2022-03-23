package servicediscovery

import (
	"consul-client/internal/config"
	"fmt"
	"log"
	"net/http"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

var (
	serviceName = "go-test-service"
)

type Client struct {
	Config          config.ServiceDiscoveryConfig
	ConsulApiConfig *consulapi.Config
	Client          *consulapi.Client
	Agent           *consulapi.Agent
	Service         *connect.Service
	KvStorage       *consulapi.KV
}

func NewClient(config config.ServiceDiscoveryConfig) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	consulApiConfig := consulapi.DefaultConfig()
	consulApiConfig.Address = config.Addr
	consulApiConfig.Token = config.Token

	c := &Client{
		Config:          config,
		ConsulApiConfig: consulApiConfig,
	}

	client, err := consulapi.NewClient(consulApiConfig)
	if err != nil {
		return nil, err
	}
	c.Client = client
	c.KvStorage = client.KV()
	c.Agent = client.Agent()

	// create service for local handling
	// svc, err := connect.NewService("go-test-consul-local", client)
	svc, err := connect.NewServiceWithConfig(serviceName, connect.Config{
		Client: c.Client,
	})
	if err != nil {
		return nil, err
	}
	c.Service = svc

	go func() {
		for {
			// fmt.Println("service ready:", c.Service.Ready())
			c.checkService()
			time.Sleep(2 * time.Second)
		}
	}()

	return c, nil
}

func (c *Client) Start(errChan chan error) {
	// c.storeData()
	// for {
	// 	c.listData()
	// 	time.Sleep(2 * time.Second)
	// }

	// c.deleteData()
	// c.getMembers()

	// c.tokenList()
	// c.test()

	c.registerService()

	c.getServices()

	// c.deregisterService()

	// c.getServices()

	fmt.Println("stop")
}

func (c *Client) ServiceListenHttp(errChan chan error) {
	server := &http.Server{
		Addr:      c.Config.HttpListenAddr,
		TLSConfig: c.Service.ServerTLSConfig(),
	}

	log.Println("Start listening http on:", c.Config.HttpListenAddr)
	errChan <- server.ListenAndServeTLS("", "")
}
