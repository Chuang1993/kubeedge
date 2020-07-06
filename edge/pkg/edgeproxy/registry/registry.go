package registry

import (
	"k8s.io/apiserver/pkg/registry/rest"
	"sync"
)

var (
	getters  map[string]rest.Getter
	listers  map[string]rest.Lister
	watchers map[string]rest.Watcher
	once     sync.Once
)

func init() {
	once.Do(func() {
		getters = make(map[string]rest.Getter)
		listers = make(map[string]rest.Lister)
		watchers = make(map[string]rest.Watcher)
	})
}

func RegistryGetter(key string, getter rest.Getter) {
	getters[key] = getter
}

func RegistryLister(key string, lister rest.Lister) {
	listers[key] = lister
}

func RegistryWatcher(key string, watcher rest.Watcher) {
	watchers[key] = watcher
}

func GetGetter(key string) rest.Getter {
	return getters[key]
}

func GetLister(key string) rest.Lister {
	return listers[key]
}
func GetWatcher(key string) rest.Watcher {
	return watchers[key]
}
