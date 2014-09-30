package client

import (
	"encoding/gob"
	"fmt"
	"github.com/cloning/go-discover/common"
	"net"
)

type Client struct {
	ServiceName string
	ServiceUrl  string
	conn        net.Conn
	closed      bool
	registry    *common.RegistrySync
}

func NewClient(name, url string) *Client {
	return &Client{name, url, nil, false, common.NewRegistrySync(nil)}
}

func (this *Client) Start() {
	this.connect()
	this.register()
	this.listen()

}

func (this *Client) Close() {

	// Indicate that connection was intentionally closed
	this.closed = true

	this.conn.Close()
}

func (this *Client) Get(serviceName string) []string {
	return this.registry.Registry[serviceName]
}

func (this *Client) listen() {
	for {
		decoder := gob.NewDecoder(this.conn)
		err := decoder.Decode(this.registry)

		if err != nil {
			// Check if connection was intentionally closed
			if this.closed {
				break
			}
			fmt.Println("Err caused client to close", err)
			break
		}
	}
}

func (this *Client) connect() {
	var err error

	this.conn, err = net.Dial("tcp", "localhost:1337")

	if err != nil {
		panic(err)
	}
}

func (this *Client) register() {
	//TODO: Replace with gob encode
	command := common.CreateRegisterCommand(this.ServiceName, this.ServiceUrl)
	fmt.Fprintf(this.conn, command+"\n")
}
