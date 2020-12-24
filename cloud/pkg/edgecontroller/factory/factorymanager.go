package factory

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"sync"
)

type InformerFactoryType string

const (
	// normal informer will be create from this factory
	Normal    InformerFactoryType = "Normal"
	Customize InformerFactoryType = "Customize"
)

type NewSharedInformerFactory func() SharedInformerFactory
type SharedInformerFactory interface {
	Start(stopCh <-chan struct{})
	ForResource(resource schema.GroupVersionResource) (informers.GenericInformer, error)
}

type informerFactoryManager struct {
	lock      sync.RWMutex
	factories map[string]SharedInformerFactory
}

func (ifm *informerFactoryManager) Load(key string) (SharedInformerFactory, bool) {
	ifm.lock.RLock()
	defer ifm.lock.Unlock()
	factory, ok := ifm.factories[key]
	return factory, ok
}

func (ifm *informerFactoryManager) Store(key string, factory SharedInformerFactory) {
	ifm.lock.Lock()
	defer ifm.lock.Unlock()
	ifm.factories[key] = factory
}

func (ifm *informerFactoryManager) LoadOrCreate(key string, newFactory NewSharedInformerFactory) SharedInformerFactory {
	factory, ok := ifm.Load(key)
	if !ok {
		factory = newFactory()
		ifm.Store(key, factory)
	}
	return factory
}


