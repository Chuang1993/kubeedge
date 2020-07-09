package watch

import "k8s.io/apimachinery/pkg/watch"

type IncomingWatcher interface {
	watch.Interface
	InComingChan() chan<- watch.Event
}

func Incoming(watcher watch.Interface) IncomingWatcher {
	if watcher == nil {
		watcher = &fakeWatcher{}
	}
	ic := &incoming{
		incoming: make(chan watch.Event),
		watcher:  watcher,
		result:   make(chan watch.Event),
	}
	go ic.loop()
	return ic
}

type incoming struct {
	incoming chan watch.Event
	watcher  watch.Interface
	result   chan watch.Event
	stopCh   chan struct{}
}

func (ic *incoming) Stop() {
	ic.stopCh <- struct{}{}
	close(ic.incoming)
	close(ic.result)
	ic.watcher.Stop()
}

func (ic *incoming) InComingChan() chan<- watch.Event {
	return ic.incoming
}

func (ic *incoming) ResultChan() <-chan watch.Event {
	return ic.result
}

func (ic *incoming) loop() {
	for {
		select {
		case event := <-ic.incoming:
			ic.result <- event
		case event := <-ic.watcher.ResultChan():
			ic.result <- event
		case <-ic.stopCh:
			return
		}
	}
}
