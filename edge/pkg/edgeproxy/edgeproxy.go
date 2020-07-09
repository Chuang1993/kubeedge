package edgeproxy

import (
	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/config"
	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/server"
	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/watch"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
)

const EdgeProxyModuleName = "edgeProxy"

type proxy struct {
	enable bool
}

func Register(proxy *v1alpha1.EdgeProxy) {
	config.InitConfigure(proxy)

	core.Register(newPorxy(proxy.Enable))
}

func newPorxy(enable bool) *proxy {
	return &proxy{enable: enable}
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
	go ps.Run()
	p.receive()
}

func (p *proxy) receive() {
	// watch请求更新的response
	go func() {
		for {
			select {
			case <-beehiveContext.Done():
				klog.Warning("stop")
				return
			default:
			}
			msg, err := beehiveContext.Receive(p.Name())
			if err != nil {
				klog.Warningf("%v", err)
				continue
			}
			watch.SendMsg(msg)
		}
	}()
}
