package server

import (
	"github.com/cloning/go-discover/common"
	"net"
	"sync"
)

const (
	PORT = ":1337"
)

type RegistrationChannelItem struct {
	command    common.RegisterCommand
	connection *Connection
}

type Server struct {
	registrationChannel     chan *RegistrationChannelItem
	closedConnectionChannel chan *Connection
	registry                *Registry
	connections             []*Connection
	running                 bool
	mutex                   *sync.Mutex
	listener                net.Listener
}

func NewServer() *Server {
	return &Server{
		make(chan *RegistrationChannelItem),
		make(chan *Connection),
		NewRegistry(),
		make([]*Connection, 0),
		false,
		&sync.Mutex{},
		nil,
	}
}

func (this *Server) Start() {
	this.startListening()
	go this.acceptRegistrations()
	go this.acceptClosedConnections()
	this.acceptConnections()
}

func (this *Server) Stop() {
	this.running = false
	for i := len(this.connections) - 1; i >= 0; i-- {
		this.connections[i].close()
	}
	this.listener.Close()
}

func (this *Server) startListening() {
	var err error
	this.listener, err = net.Listen("tcp", PORT)

	if err != nil {
		panic(err)
	}
	this.running = true
}

func (this *Server) acceptConnections() {
	for {

		conn, err := this.listener.Accept()

		if err != nil {
			if this.running == false {
				// Server has been shut down, causing the problem
				break
			}
		}

		this.handleConnection(conn)
	}
}

func (this *Server) handleConnection(conn net.Conn) {
	connection := NewConnection(conn, this)
	this.addConnection(connection)
	go connection.run()
}

func (this *Server) acceptRegistrations() {
	for {
		item := <-this.registrationChannel
		this.registry.add(item.command, item.connection)
		this.syncRegistry()
	}
}

func (this *Server) acceptClosedConnections() {
	for {
		connToClose := <-this.closedConnectionChannel
		this.removeConnection(connToClose)
		this.registry.remove(connToClose)
		this.syncRegistry()
	}
}

func (this *Server) syncRegistry() {
	syncData := this.registry.get()
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, conn := range this.connections {
		conn.sendRegistry(syncData)
	}
}

func (this *Server) addConnection(conn *Connection) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.connections = append(this.connections, conn)
}

func (this *Server) removeConnection(connToClose *Connection) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for i, conn := range this.connections {
		if conn == connToClose {
			this.connections = append(this.connections[:i], this.connections[i+1:]...)
			break
		}
	}
}
