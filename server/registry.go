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

	for k, v := range this.register {
		for i, item := range v {
			if item.connection == connection {
				this.register[k] = append(v[:i], v[i+1:]...)
				return
			}
		}
	}
}

func (this *Registry) get() *common.RegistrySync {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	serviceRegister := make(map[string][]string)
	for k, v := range this.register {
		endpoints := make([]string, len(v))
		for i, item := range v {
			endpoints[i] = item.serviceUrl
		}
		serviceRegister[k] = endpoints
	}
	return common.NewRegistrySync(serviceRegister)
}

func NewRegistry() *Registry {
	return &Registry{
		make(map[string][]*RegistryItem),
		&sync.Mutex{},
	}
}
