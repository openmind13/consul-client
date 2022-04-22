package consul

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type Client struct {
	mu     sync.RWMutex
	config Config

	serviceID      string
	serviceVersion string
	checkId        string
	nodeName       string

	consulApiConfig *consulapi.Config
	client          *consulapi.Client
	agent           *consulapi.Agent
	catalog         *consulapi.Catalog
	kvStorage       *consulapi.KV
}

func NewClient(config Config, serviceVersion string) (*Client, error) {
	consulApiConfig := consulapi.DefaultConfig()
	consulApiConfig.Address = config.Address
	consulApiConfig.Token = config.Token

	c := &Client{
		serviceID:       config.ServiceID,
		checkId:         config.CheckName + "-" + uuid.NewString(),
		config:          config,
		consulApiConfig: consulApiConfig,
		serviceVersion:  serviceVersion,
	}

	client, err := consulapi.NewClient(consulApiConfig)
	if err != nil {
		return nil, fmt.Errorf("consul error: %s", err.Error())
	}
	c.client = client
	c.kvStorage = client.KV()
	c.agent = client.Agent()
	c.catalog = client.Catalog()

	if err := c.registerService(); err != nil {
		return nil, fmt.Errorf("consul error: %s", err.Error())
	}

	nodeName, err := c.agent.NodeName()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to get node name from consul")
	} else {
		c.nodeName = nodeName
		logrus.Debug("node name = ", nodeName)
	}

	go c.StartConfigWatcher()

	return c, nil
}

func (c *Client) Start(errChan chan error) {
	go c.updateServiceLoop()

	go func() {
		for {
			// servicesData, meta, err := c.catalog.Services(&consulapi.QueryOptions{})
			// if err != nil {
			// 	logrus.Error(err)
			// } else {
			// 	logrus.Info("time:", meta.RequestTime)
			// 	for key, svcData := range servicesData {
			// 		fmt.Println(key, svcData)
			// 	}
			// }

			// nodes, _, err := c.catalog.Nodes(&consulapi.QueryOptions{})
			// if err != nil {
			// 	logrus.Error(err)
			// } else {
			// 	fmt.Println()
			// 	for _, node := range nodes {
			// 		fmt.Println(node.Node, node.Meta)
			// 	}
			// }

			// meta, err := c.agent.Host()
			// if err != nil {
			// 	logrus.Error(err)
			// } else {
			// 	fmt.Println()
			// 	for key, value := range meta {
			// 		fmt.Println(key, value)
			// 	}
			// }

			nodeServices, _, err := c.catalog.NodeServiceList("consul-server-0", &consulapi.QueryOptions{})
			if err != nil {
				logrus.Error(err)
			} else {
				fmt.Println()
				fmt.Println(nodeServices.Node.Node)
				for _, svc := range nodeServices.Services {
					fmt.Println(svc.Service)
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	select {}
}

func (c *Client) updateServiceLoop() {
	for {
		if err := c.agent.UpdateTTL(c.checkId, "", "pass"); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to update consul service")
		}
		time.Sleep(c.GetConfig().UpdateCheckPeriod)
	}
}

func (c *Client) IsReady() bool {
	return true
}

func (c *Client) registerService() error {
	checks := consulapi.AgentServiceChecks{}

	httpCalcServiceCheck := &consulapi.AgentServiceCheck{
		Name:    c.GetConfig().CheckName,
		CheckID: c.checkId,

		TTL: c.config.CheckTTL.String(),
	}
	checks = append(checks, httpCalcServiceCheck)

	// redisCheck := &consulapi.AgentServiceCheck{
	// 	CheckID: REDIS_CHECK_ID,
	// 	Name:    "redis-check-name",
	// 	TTL:     "5s",
	// }
	// checks = append(checks, redisCheck)

	serviceRegister := consulapi.AgentServiceRegistration{
		ID:   c.serviceID,
		Name: c.GetConfig().ServiceName,
		Meta: map[string]string{
			"version": c.serviceVersion,
		},
		Tags:   c.config.Tags,
		Checks: checks,
		TaggedAddresses: map[string]consulapi.ServiceAddress{
			"localhost": {
				Address: "localhost",
				Port:    8080,
			},
		},
	}

	if err := c.agent.ServiceRegister(&serviceRegister); err != nil {
		return err
	}
	return nil
}

// var (
// 	CATALOG_NODE_NAME = "consul-server-0"
// 	CATALOG_ID        = ""
// )

// func (c *Client) RegisterServiceCatalog() error {
// 	// svc, _, err := c.agent.Service(c.serviceId, &consulapi.QueryOptions{})
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	catalogRegistration := consulapi.CatalogRegistration{
// 		ID:   CATALOG_ID,
// 		Node: CATALOG_NODE_NAME,
// 		NodeMeta: map[string]string{
// 			"field": "some_meta",
// 		},
// 		TaggedAddresses: map[string]string{
// 			"some": "test",
// 		},
// 		SkipNodeUpdate: true,
// 		// Service:        svc,
// 		Address: "localhost",
// 	}
// 	_, err := c.catalog.Register(&catalogRegistration, &consulapi.WriteOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (c *Client) DeregisterServiceCatalog() error {
// 	catalogDeregister := consulapi.CatalogDeregistration{
// 		Node: CATALOG_NODE_NAME,
// 	}
// 	if _, err := c.catalog.Deregister(&catalogDeregister, &consulapi.WriteOptions{}); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c *Client) StopService() error {
	return c.client.Agent().ServiceDeregister(c.serviceID)
}

func (c *Client) GetServices() ([]consulapi.AgentServiceChecksInfo, error) {
	svcMap, err := c.agent.ServicesWithFilterOpts("", &consulapi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	services := []consulapi.AgentServiceChecksInfo{}

	for _, service := range svcMap {
		_, meta, err := c.agent.AgentHealthServiceByIDOpts(service.ID, &consulapi.QueryOptions{})
		// status, meta, err := c.agent.AgentHealthServiceByNameOpts(svcName, &consulapi.QueryOptions{})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to check services in consul")
			continue
		}
		if meta != nil {
			services = append(services, *meta)
		}
	}

	return services, nil
}

func (c *Client) DeregisterServiceWithID(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

func (c *Client) SetConfig(conf Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = conf
}

func (c *Client) GetConfig() Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

func (c *Client) GetNodeName() string {
	return c.nodeName
}
