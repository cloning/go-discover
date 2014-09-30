package server

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/cloning/go-discover/common"
	"io"
	"net"
)

type Connection struct {
	conn   net.Conn
	server *Server
}

func NewConnection(conn net.Conn, server *Server) *Connection {
	return &Connection{conn, server}
}

func (this *Connection) close() {
	this.conn.Close()
}

func (this *Connection) run() {
	for {
		reader := bufio.NewReader(this.conn)

		msg, err := reader.ReadString('\n')

		// Client disconnected
		if err == io.EOF {
			this.server.closedConnectionChannel <- this
			break
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		this.handleMsg(msg)
	}
}

func (this *Connection) sendRegistry(registry *common.RegistrySync) {
	enc := gob.NewEncoder(this.conn)
	enc.Encode(registry)
}

func (this *Connection) handleMsg(msg string) {

	command := common.ParseCommand(msg)

	if w, ok := command.(common.RegisterCommand); ok {
		this.server.registrationChannel <- &RegistrationChannelItem{w, this}
	}
}
