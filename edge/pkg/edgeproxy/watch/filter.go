package watch

import "k8s.io/apimachinery/pkg/watch"

func AlwaysTrue(in watch.Event) (watch.Event, bool) {
	return in, true
}
