package server

import (
	"fmt"
	"github.com/cloning/go-discover/common"
	"net"
	"sync"
	"time"
)

const (
	REAPER_THREAD_SLEEP_MS = 100
	PORT                   = ":1337"
)

type Server struct {
	registrationChannel chan common.RegisterCommand
	registry            *common.Registry
	connections         []Connection
	running             bool
	mutex               *sync.Mutex
	listener            net.Listener
}

func (this *Server) Start() {
	this.initialize()
	this.startListening()
	go this.reap()
	go this.acceptRegistrations()
	this.acceptConnections()
}

func (this *Server) initialize() {
	this.registry = common.CreateRegistry()
	this.registrationChannel = make(chan common.RegisterCommand)
	this.mutex = &sync.Mutex{}
	this.running = true
	this.connections = make([]Connection, 0)
}

func (this *Server) startListening() {
	var err error
	this.listener, err = net.Listen("tcp", PORT)

	if err != nil {
		panic(err)
	}
}

func (this *Server) acceptConnections() {
	for {
		if this.running == false {
			break
		}

		conn, err := this.listener.Accept()

		if err != nil {
			fmt.Println(err)
		}

		this.handleConnection(conn)
	}
}

func (this *Server) handleConnection(conn net.Conn) {
	connection := Connection{conn, this, false, this.registrationChannel}
	this.mutex.Lock()
	this.connections = append(this.connections, connection)
	this.mutex.Unlock()
	go connection.run()
}

func (this *Server) acceptRegistrations() {
	for {
		this.registry.Add(<-this.registrationChannel)
	}
}

func (this *Server) reap() {
	for {
		if this.running == false {
			break
		}

		this.mutex.Lock()

		for i := len(this.connections) - 1; i >= 0; i-- {
			conn := this.connections[i]
			if !conn.Closed {
				this.connections = append(this.connections[:i], this.connections[i+1:]...)
			}
		}

		this.mutex.Unlock()

		time.Sleep(REAPER_THREAD_SLEEP_MS * time.Millisecond)
	}
}
