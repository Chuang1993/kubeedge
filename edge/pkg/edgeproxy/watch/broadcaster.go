package watch

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/core/model"
)

const AddWatcherOperation = "addwatcher"

func NewMsgBrocaster() *MsgBrocaster {
	brocaster := &MsgBrocaster{
		msgQueue: make(chan model.Message, 1024),
		watchers: sync.Map{},
		cache:    newCache(),
	}
	go brocaster.loop()
	return brocaster
}

type MsgBrocaster struct {
	msgQueue chan model.Message
	watchers sync.Map
	cache    *cache
}

type WatcherAddFunc func()

type addWatcher struct {
	AddFunc WatcherAddFunc
}

func (mb *MsgBrocaster) Brocast(msg model.Message) {
	mb.msgQueue <- msg
}

func (mb *MsgBrocaster) loop() {
	for msg := range mb.msgQueue {
		if msg.Router.Operation == AddWatcherOperation {
			aw := msg.Content.(addWatcher)
			aw.AddFunc()
			continue
		}
		mb.cache.Add(msg)
		mb.watchers.Range(func(key, value interface{}) bool {
			watcher := value.(IncomingWatcher)
			watcher.InComingChan() <- ConvertMsgToEvent(msg)
			return true
		})
	}
}

func (mb *MsgBrocaster) AddWatcher(key string, watcher IncomingWatcher, timeout time.Duration, resourceVersion uint64) {
	msg := model.Message{
		Header: model.MessageHeader{},
		Router: model.MessageRoute{Operation: AddWatcherOperation},
		Content: addWatcher{
			AddFunc: func() {
				msgs, err := mb.cache.GetCacheMsg(resourceVersion)
				if err != nil {
					watcher.InComingChan() <- watch.Event{
						Type:   watch.Error,
						Object: nil,
					}
					klog.Error("Err: %v", err)
					return
				}
				for i := range msgs {
					watcher.InComingChan() <- ConvertMsgToEvent(msgs[i])
				}
				mb.watchers.Store(key, watcher)
				go func(td time.Duration) {
					tc := time.After(td)
					<-tc
					klog.Warning("watcher timeout, remove watcher from brocaster")
					value, _ := mb.watchers.Load(key)
					watcher := value.(IncomingWatcher)
					watcher.Stop()
					mb.RemoveWatcher(key)

				}(timeout)
			}},
	}
	mb.Brocast(msg)
}

func (mb *MsgBrocaster) RemoveWatcher(key string) {
	mb.watchers.Delete(key)
}
