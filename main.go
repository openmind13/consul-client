package main

import (
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

// const addr = "0.0.0.0:8600"
const addr = "172.17.0.2:8600"

func main() {
	log.Println("consul client")

	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("Failed to create client", err)
	}

	svc, err := connect.NewService("counting", client)
	if err != nil {
		log.Fatal(err)
	}
	defer svc.Close()

	log.Println("connected")

	agent := client.Agent()

	registration := &api.AgentServiceRegistration{
		ID:      "counting",
		Name:    "counting",
		Port:    9001,
		Address: "0.0.0.0",
	}

	err = agent.ServiceRegister(registration)
	if err != nil {
		log.Fatal("register error ", err)
	}

	// members, err := client.Agent().Members(true)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(members)

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// go func() {
	// 	select {
	// 	case <-ctx.Done():
	// 	case <-time.After(2 * time.Second):
	// 		cancel()
	// 		log.Fatal("context exceeded")
	// 	}
	// }()

	values, meta, err := client.KV().Get("test", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(values, meta)

	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: svc.ServerTLSConfig(),
	}

	log.Println("server started")
	server.ListenAndServe()
}
