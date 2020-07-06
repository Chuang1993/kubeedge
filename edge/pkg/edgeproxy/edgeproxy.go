package edgeproxy

import (
	"github.com/kubeedge/beehive/pkg/common/util"
	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/config"
	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/server"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/client"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
	"k8s.io/klog"
)

const EdgeProxyModuleName = "edgeProxy"

type proxy struct {
	enable bool
	// for get operation
	metaClient client.CoreInterface
}

func Register(proxy *v1alpha1.EdgeProxy) {
	config.InitConfigure(proxy)

	core.Register(newPorxy(proxy.Enable))
}

func newPorxy(enable bool) *proxy {
	metaClient := client.New()
	return &proxy{enable: enable, metaClient: metaClient}
}
func (p *proxy) Name() string {
	return EdgeProxyModuleName
}
func (p *proxy) Group() string {
	return modules.ProxyGroup
}

func (p *proxy) Enable() bool {
	return p.enable
}
func (p *proxy) Start() {
	// TODO start http server with Config.ListenPort
	// new proxyServer，异步run
	// 是否考虑使用stopch停止服务
	ps := server.NewProxyServer()
	ps.Run()
	// 是否需要写到proxy这儿？这儿更加倾向于从metaManager过来的消息分发处理
	p.runProxy()
}

func (p *proxy) runProxy() {
	// watch请求更新的response
	go func() {
		for {
			select {
			case <-beehiveContext.Done():
				klog.Warning("stop")
				return
			default:
			}
			if msg, err := beehiveContext.Receive(p.Name()); err != nil {
				p.process(msg)
			}
		}
	}()
}
func (p *proxy) process(msg model.Message) error {
	//TODO 根据不同的resouce来分派消息
	// 考虑实现不同的resouce的dispatch内容，根据msg的route的operation来判断对于client的event类型
	_, msgType, _, err := util.ParseResourceEdge(msg.GetResource(), msg.GetOperation())
	if err != nil {
		klog.Error("error")
		return err
	}
	switch msgType {
	case "pod":
	// send msg to pod watch request
	case "service":
	//send msg to svc watch request
	case "endpoints":
		// send msg to ep watch request
	}
	return nil
}
