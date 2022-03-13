package servicediscovery

import (
	"consul-client/config"
	"fmt"
	"log"
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

const (
	HTTP_ADDR = "0.0.0.0:8080"
)

type Client struct {
	Config    *consulapi.Config
	Client    *consulapi.Client
	Agent     *consulapi.Agent
	Service   *connect.Service
	KvStorage *consulapi.KV
}

func NewClient(config config.ServiceDiscoveryConfig) (*Client, error) {
	sdConfig := consulapi.DefaultConfig()
	sdConfig.Address = config.Addr
	sdConfig.Token = config.Token

	c := &Client{
		Config: sdConfig,
	}

	client, err := consulapi.NewClient(sdConfig)
	if err != nil {
		return nil, err
	}
	c.Client = client
	c.KvStorage = client.KV()
	c.Agent = client.Agent()

	svc, err := connect.NewService("go-test-consul", client)
	if err != nil {
		return nil, err
	}
	c.Service = svc

	return c, nil
}

func (c *Client) Start(errChan chan error) {
	// c.storeData()
	// for {
	// 	c.listData()
	// 	time.Sleep(2 * time.Second)
	// }

	c.deleteData()
}

func (c *Client) storeData() {
	pair := consulapi.KVPair{
		Key:   "test_key",
		Value: []byte("test_value"),
	}
	meta, err := c.KvStorage.Put(&pair, &consulapi.WriteOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("writed in:", meta.RequestTime)
}

func (c *Client) listData() {
	pairs, meta, err := c.KvStorage.List("/", &consulapi.QueryOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if pairs == nil || meta == nil {
		log.Fatal("pairs or meta is nil")
	}
	for _, pair := range pairs {
		fmt.Println(pair.Key, string(pair.Value))
	}
}

func (c *Client) deleteData() {
	meta, err := c.KvStorage.Delete("test_key", &consulapi.WriteOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted in:", meta.RequestTime)
}

func (c *Client) getMembers() {
	members, err := c.Agent.Members(false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("members:")
	for _, member := range members {
		fmt.Println(member.Name)
	}
	fmt.Println()
}

func (c *Client) getMetrics() {
	metrics, err := c.Agent.Metrics()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("metrics:")
	fmt.Printf("%+v\n", metrics)
	fmt.Println()
}

func (c *Client) Listen(errChan chan error) {
	server := &http.Server{
		Addr:      HTTP_ADDR,
		TLSConfig: c.Service.ServerTLSConfig(),
	}
	log.Println("Start listening http on:", HTTP_ADDR)
	errChan <- server.ListenAndServeTLS("", "")
}
