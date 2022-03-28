package servicediscovery

import (
	"consul-client/internal/config"
	"log"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

const (
	// SERVICE_ID   = "go-test-service-dev"
	SERVICE_NAME = "go-test-service-dev"
	CHECK_ID     = "check-go-id"
	CHECK_NAME   = "go-check-name"
)

type Client struct {
	config          config.Config
	consulApiConfig *consulapi.Config
	client          *consulapi.Client
	agent           *consulapi.Agent
	service         *connect.Service
	kvStorage       *consulapi.KV
}

func NewClient(config config.Config) (*Client, error) {

	consulApiConfig := consulapi.DefaultConfig()
	consulApiConfig.Address = config.ConsulAddress
	consulApiConfig.Token = config.Token

	c := &Client{
		config:          config,
		consulApiConfig: consulApiConfig,
	}

	client, err := consulapi.NewClient(consulApiConfig)
	if err != nil {
		return nil, err
	}
	c.client = client
	c.kvStorage = client.KV()
	c.agent = client.Agent()

	// create service for local handling
	// svc, err := connect.NewService(serviceName, client)
	// // svc, err := connect.NewServiceWithConfig(serviceName, connect.Config{
	// // 	Client: c.Client,
	// // })
	// if err != nil {
	// 	return nil, err
	// }
	// c.Service = svc

	// go func() {
	// 	for {
	// 		// fmt.Println("service ready:", c.Service.Ready())
	// 		// c.checkService()
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	return c, nil
}

func (c *Client) StartService(errChan chan error) {

	// go func() {
	// 	for {
	// 		// c.checkService()
	// 		// c.getServices()

	// 		if err := c.agent.UpdateTTL(CHECK_ID, "output_test", "passing"); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		time.Sleep(time.Second)
	// 	}
	// }()

}

func (c *Client) StartServiceUpdater() {
	for {
		if err := c.agent.UpdateTTL(CHECK_ID, "go test check passed", "passing"); err != nil {
			// log.Fatal(err)
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
}

func (c *Client) RegisterService() error {
	checks := consulapi.AgentServiceChecks{}

	httpCheck := consulapi.AgentServiceCheck{
		CheckID: CHECK_ID,
		// Name:    CHECK_NAME,
		TTL: "5s",
	}

	checks = append(checks, &httpCheck)

	serviceRegister := consulapi.AgentServiceRegistration{
		ID:   SERVICE_NAME,
		Name: SERVICE_NAME,
		Meta: map[string]string{
			"test_key": "test_value",
		},
		Tags:   []string{"dev", "test-client", "go"},
		Checks: checks,
	}
	if err := c.client.Agent().ServiceRegister(&serviceRegister); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeregisterService() error {
	return c.client.Agent().ServiceDeregister(SERVICE_NAME)
}

// func (c *Client) ServiceListenHttp(errChan chan error) {
// 	server := &http.Server{
// 		Addr:      c.config.HttpListenAddress,
// 		TLSConfig: c.service.ServerTLSConfig(),
// 	}
// 	log.Println("Start listening http on:", c.config.HttpListenAddress)
// 	errChan <- server.ListenAndServeTLS("", "")
// }
