package common

import (
	"fmt"
	"sync"
)

type Registry struct {
	register map[string][]string
	mutex    *sync.Mutex
}

func (this *Registry) Add(registration RegisterCommand) {
	this.mutex.Lock()
	this.register[registration.ServiceName] = append(
		this.register[registration.ServiceName],
		registration.ServiceUrl)
	fmt.Println(this.register)
	this.mutex.Unlock()
}

func CreateRegistry() *Registry {
	return &Registry{
		make(map[string][]string),
		&sync.Mutex{},
	}
}
