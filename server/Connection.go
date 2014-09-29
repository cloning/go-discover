package server

import (
	"bufio"
	"fmt"
	"github.com/cloning/go-discover/common"
	"io"
	"net"
)

type Connection struct {
	conn                net.Conn
	server              *Server
	Closed              bool
	registrationChannel chan common.RegisterCommand
}

func (this *Connection) run() {
	for {
		reader := bufio.NewReader(this.conn)

		msg, err := reader.ReadString('\n')

		// Client disconnected
		if err == io.EOF {
			this.Closed = true
			fmt.Println("Connection is the closed")
			break
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		this.handleMsg(msg)
	}
}

func (this *Connection) handleMsg(msg string) {

	command := common.ParseCommand(msg)

	if w, ok := command.(common.RegisterCommand); ok {
		this.registrationChannel <- w
	}
}
