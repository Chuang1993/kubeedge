package config

import (
	"github.com/kubeedge/kubeedge/cloud/pkg/edgecontroller/factory"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"sync"

	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
)

var Config Configure
var once sync.Once

type Configure struct {
	v1alpha1.EdgeController
	KubeAPIConfig  v1alpha1.KubeAPIConfig
	NodeName       string
	EdgeSiteEnable bool
}

func InitConfigure(ec *v1alpha1.EdgeController, kubeAPIConfig *v1alpha1.KubeAPIConfig, nodeName string, edgesite bool) {
	once.Do(func() {
		Config = Configure{
			EdgeController: *ec,
			KubeAPIConfig:  *kubeAPIConfig,
			NodeName:       nodeName,
			EdgeSiteEnable: edgesite,
		}
	})
}

type DownStreamResource struct {
	ResourceEventHandler cache.ResourceEventHandler
	InformerFactory      informers.SharedInformerFactory
}

func GetCommonResourceEventHandler() cache.ResourceEventHandler {

	events := make(chan watch.Event)
	return factory.NewCommonResourceEventHandler(events)
}

var DownStreamResources = []DownStreamResource{
	DownStreamResource{

	},
	DownStreamResource{

	},
}
