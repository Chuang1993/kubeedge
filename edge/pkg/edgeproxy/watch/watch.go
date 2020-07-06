package watch

import (
	"github.com/kubeedge/beehive/pkg/core/model"
	"k8s.io/apimachinery/pkg/runtime"
	"sync"
)

// 直接在sqlite内构建筛选，不是特别合适，考虑自己实现。
type FilterFunc func(runtime.Object) bool

type WatchChan struct {
	// 容量是否为1024？
	inChan     chan model.Message
	filterFunc FilterFunc
	// 是否完成缓存内信息发送，直接接收广播消息
	ready bool
}

type cacheElement struct {
	value    []runtime.Object
	startIdx int
	capacity int
}

func (elem *cacheElement) Add(obj runtime.Object) {

}

type WatchMsgMgr struct {
	// 完成广播，让同一个runtime.Object进入到所有的watcher接口内
	watchChan map[string][]*WatchChan
	wcLock    sync.RWMutex
	// 缓存部分数据，缓存一部分数据，避免list和watch之间间隔，大批量数据被忽略。如果失败，返回异常，让client端重新尝试
	cache map[string]*cacheElement
	cLock sync.RWMutex
}

func (wmm *WatchMsgMgr) SendMsg(msg model.Message) {
	// TODO add cache& delete experied msg,keep 100 msg only
	// TODO 完成消息广播
	wmm.wcLock.RLock()
	defer wmm.wcLock.RUnlock()

}
func (wmm *WatchMsgMgr) AddWatcher(resource string, wc *WatchChan) {
	wmm.wcLock.Lock()
	defer wmm.wcLock.Unlock()
	wcs := wmm.watchChan[resource]
	wcs = append(wcs, wc)
	wmm.watchChan[resource] = wcs
	// sendCacheMsg
}
