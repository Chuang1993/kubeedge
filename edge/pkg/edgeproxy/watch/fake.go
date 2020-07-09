package watch

import "k8s.io/apimachinery/pkg/watch"

type fakeWatcher struct {
	result chan watch.Event
}

func (fw *fakeWatcher) Stop() {
	close(fw.result)
}
func (fw *fakeWatcher) ResultChan() <-chan watch.Event {
	return fw.result
}
