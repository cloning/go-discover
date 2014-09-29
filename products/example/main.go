package main

import (
	"github.com/cloning/go-discover/client"
	"github.com/cloning/go-discover/server"
)

func main() {
	server := &server.Server{}
	go server.Start()

	c1 := client.NewClient("users", "localhost:8080")
	go c1.Start()

	c2 := client.NewClient("users", "localhost:8090")
	c2.Start()
}
