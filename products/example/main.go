package main

import (
	"fmt"
	"github.com/cloning/go-discover/client"
	"github.com/cloning/go-discover/server"
	"time"
)

func main() {
	server := server.NewServer()
	go server.Start()

	c1 := client.NewClient("users", "localhost:8080")
	go c1.Start()

	time.Sleep(1 * time.Second)
	fmt.Println("Getting service from c1", c1.Get("users"))
	time.Sleep(1 * time.Second)

	c2 := client.NewClient("users", "localhost:8090")
	go c2.Start()

	time.Sleep(1 * time.Second)
	fmt.Println("Getting service from c1", c1.Get("users"))
	fmt.Println("Getting service from c2", c2.Get("users"))
	time.Sleep(1 * time.Second)

	time.Sleep(1 * time.Second)
	c1.Close()

	time.Sleep(1 * time.Second)
	fmt.Println("Getting service from c2", c2.Get("users"))
	time.Sleep(1 * time.Second)
	c2.Close()
	time.Sleep(1 * time.Second)
	server.Stop()
	time.Sleep(1 * time.Second)
}
