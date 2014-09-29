package client

import (
	"bufio"
	"fmt"
	"github.com/cloning/go-discover/common"
	"net"
)

type Client struct {
	ServiceName string
	ServiceUrl  string
	conn        net.Conn
	closed      bool
}

func NewClient(name, url string) *Client {
	return &Client{name, url, nil, false}
}

func (this *Client) Start() {
	this.connect()
	this.register()
	this.listen()

}

func (this *Client) listen() {
	for {
		// TODO: Handle internal registry updates from server here
		status, err := bufio.NewReader(this.conn).ReadString('\n')

		if err != nil {
			panic(err)
		}

		fmt.Println(status)
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
	command := common.CreateRegisterCommand(this.ServiceName, this.ServiceUrl)
	fmt.Fprintf(this.conn, command+"\n")
}
