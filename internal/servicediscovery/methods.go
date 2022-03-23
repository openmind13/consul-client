package servicediscovery

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

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

func (c *Client) getServices() {
	servicesMap, err := c.Agent.Services()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("services:")
	for name, service := range servicesMap {
		fmt.Println(name, service.Address, service.Datacenter, service.SocketPath)
	}
}

// func (c *Client) test() {
// 	str, data, err := c.Agent.AgentHealthServiceByID("go-test-consul")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(str, data)
// }

func (c *Client) registerService() {
	check := consulapi.AgentServiceCheck{
		// Interval: "2s",
		TTL: "2s",
	}
	meta := map[string]string{}
	meta["test"] = "hello"
	serviceReg := consulapi.AgentServiceRegistration{
		Kind: consulapi.ServiceKindTypical,
		// Namespace: "test",
		Check: &check,
		ID:    serviceName,
		Name:  "test",
		Port:  8080,
		// Address: ,
		Meta: meta,
	}
	if err := c.Client.Agent().ServiceRegister(&serviceReg); err != nil {
		log.Fatal(err)
	}
}

func (c *Client) deregisterService() {
	if err := c.Client.Agent().ServiceDeregister(serviceName); err != nil {
		log.Println("failed to deregister service", err)
	}
}

func (c *Client) checks() {
	checks, err := c.Client.Agent().Checks()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(checks)
}

func (c *Client) checkService() {
	status, info, err := c.Client.Agent().AgentHealthServiceByID(serviceName)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("service name:", serviceName, "status:", status)
	if info != nil {
		fmt.Printf("%+v\n", info.Service)
	}
}
