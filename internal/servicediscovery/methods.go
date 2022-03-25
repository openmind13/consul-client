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
	meta, err := c.kvStorage.Put(&pair, &consulapi.WriteOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("writed in:", meta.RequestTime)
}

func (c *Client) listData() {
	pairs, meta, err := c.kvStorage.List("/", &consulapi.QueryOptions{})
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
	meta, err := c.kvStorage.Delete("test_key", &consulapi.WriteOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted in:", meta.RequestTime)
}

func (c *Client) getMembers() {
	members, err := c.agent.Members(false)
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
	metrics, err := c.agent.Metrics()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("metrics:")
	fmt.Printf("%+v\n", metrics)
	fmt.Println()
}

func (c *Client) getServices() {
	servicesMap, err := c.agent.Services()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("services:")
	for name, service := range servicesMap {
		fmt.Println(name, service.Address, service.Datacenter, service.SocketPath)
	}
}

func (c *Client) printChecks() {
	checks, err := c.client.Agent().Checks()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(checks)
}

func (c *Client) checkService() {
	status, info, err := c.client.Agent().AgentHealthServiceByID(SERVICE_NAME)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("service name:", SERVICE_NAME, "status:", status)
	if info != nil {
		fmt.Printf("%+v\n", info.Service)

		// if len(info.Checks) == 0 {
		// 	fmt.Println("checks is empty")
		// } else {
		// 	for _, check := range info.Checks {
		// 		fmt.Printf("check: %+v\n", check)
		// 	}
		// }

	}

	fmt.Println()
}
