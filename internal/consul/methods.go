package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

func (c *Client) storeData() {
	pair := consulapi.KVPair{
		Key:   "test_key",
		Value: []byte("test_value"),
	}
	meta, err := c.kvStorage.Put(&pair, &consulapi.WriteOptions{})
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println("writed in:", meta.RequestTime)
}

func (c *Client) listData() {
	pairs, meta, err := c.kvStorage.List("/", &consulapi.QueryOptions{})
	if err != nil {
		logrus.Error(err)
	}
	if pairs == nil || meta == nil {
		logrus.Error("pairs or meta is nil")
	}
	for _, pair := range pairs {
		fmt.Println(pair.Key, string(pair.Value))
	}
}

func (c *Client) deleteData() {
	meta, err := c.kvStorage.Delete("test_key", &consulapi.WriteOptions{})
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println("deleted in:", meta.RequestTime)
}

func (c *Client) getMembers() {
	members, err := c.agent.Members(true)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println("members:")
	for _, member := range members {
		fmt.Println(member.Name)
		for tagName, tag := range member.Tags {
			fmt.Println(tagName, tag)
		}
	}
	fmt.Println()
}

func (c *Client) getMetrics() {
	metrics, err := c.agent.Metrics()
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println("metrics:")
	fmt.Printf("%+v\n", metrics)
	fmt.Println()
}

func (c *Client) getServices() {
	servicesMap, err := c.agent.Services()
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println("services:")
	for name, service := range servicesMap {
		fmt.Println(name, service.Address, service.Datacenter, service.SocketPath)
	}
}

func (c *Client) test() {
	str, data, err := c.agent.AgentHealthServiceByID("go-test-consul")
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(str, data)
}

func (c *Client) checks() {
	checks, err := c.client.Agent().Checks()
	if err != nil {
		logrus.Error(err)
	}
	for name, check := range checks {
		fmt.Printf("%s %+v\n", name, check)
	}
	// fmt.Println(checks)
}

// func (c *Client) checkService() {
// 	status, info, err := c.client.Agent().AgentHealthServiceByID(SERVICE_NAME)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	fmt.Println("service name:", SERVICE_NAME, "status:", status)
// 	if info != nil {
// 		// fmt.Printf("%+v\n", info.Service)

// 		// if len(info.Checks) == 0 {
// 		// 	fmt.Println("checks is empty")
// 		// } else {
// 		// 	for _, check := range info.Checks {
// 		// 		fmt.Printf("check: %+v\n", check)
// 		// 	}
// 		// }

// 	}

// 	fmt.Println()
// }
