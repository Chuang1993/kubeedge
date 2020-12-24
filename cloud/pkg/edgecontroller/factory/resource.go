package factory

import (
	beehiveModel "github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/cloud/pkg/edgecontroller/config"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

type RunInformerFunc func() chan<- watch.Event
type FilterEventFunc func(event watch.Event) bool

type ConvertEventToMsgFunc func(event watch.Event, nodeid string) *beehiveModel.Message

type SendMsgFunc func(msg *beehiveModel.Message)

/*
1. 获取listwatch的结果
2. 判断事件是否发送
3. 获取路由
4. 转换为beehive消息
5. 完成发送
*/

type DownFlow struct {
	Strategy RouteStrategy
	Informer RunInformerFunc
	Filter   FilterEventFunc
	Router   *Router
	Sender   *Sender
}

func (dr *DownFlow) Start() {
	eventChan := dr.Informer()
	for event := range eventChan {
		if dr.Filter(event) {
			nodeids := dr.Router.Route(event, dr.Strategy)
			for _, nodeid := range nodeids {
				msg := dr.Sender.Convert(event, nodeid)
				dr.Sender.Send(msg)
			}
		}
	}
}

type Sender struct {
	Convert ConvertEventToMsgFunc
	Send    SendMsgFunc
}

var scheme = runtime.NewScheme()
var stopNever = make(chan struct{})

func RunInformer(kubeClient *kubernetes.Clientset, kind schema.GroupVersionKind, resource, namespace string, resyncPeriod time.Duration, fieldSelector fields.Selector) (chan<- watch.Event, error) {
	lw := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), resource, namespace, fieldSelector)
	events := make(chan watch.Event, config.Config.Buffer.ConfigMapEvent)
	rh := NewCommonResourceEventHandler(events)
	obj, err := scheme.New(kind)
	if err != nil {
		return nil, err
	}
	si := cache.NewSharedInformer(lw, obj, resyncPeriod)
	si.AddEventHandler(rh)
	go si.Run(stopNever)
	return events, nil
}
