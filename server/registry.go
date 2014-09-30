package server

import (
	"github.com/cloning/go-discover/common"
	"sync"
)

type RegistryItem struct {
	serviceUrl string
	connection *Connection
}

type Registry struct {
	register map[string][]*RegistryItem
	mutex    *sync.Mutex
}

func (this *Registry) add(registration common.RegisterCommand, connection *Connection) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.register[registration.ServiceName] = append(
		this.register[registration.ServiceName],
		&RegistryItem{
			registration.ServiceUrl,
			connection})
}

func (this *Registry) remove(connection *Connection) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, v := range this.register {
		for i, item := range v {
			if item.connection == connection {
				v = append(v[:i], v[i+1:]...)
				return
			}
		}
	}
}

func NewRegistry() *Registry {
	return &Registry{
		make(map[string][]*RegistryItem),
		&sync.Mutex{},
	}
}
