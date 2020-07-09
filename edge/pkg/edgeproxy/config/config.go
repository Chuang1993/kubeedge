package config

import (
	"sync"

	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
)

var Config Configure
var once sync.Once

type Configure struct {
	v1alpha1.EdgeProxy
}

func InitConfigure(c *v1alpha1.EdgeProxy) {
	once.Do(func() {
		Config = Configure{EdgeProxy: *c}

	})
}
