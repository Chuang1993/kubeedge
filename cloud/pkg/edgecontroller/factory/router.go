package factory

import (
	"k8s.io/apimachinery/pkg/watch"
	"sync"
)

type RouteStrategy string

const (
	NodeStrategy RouteStrategy = "Node"
	PodStrategy  RouteStrategy = "Pod"
	ALLStrategy  RouteStrategy = "All"
)

type RouteFunc func(event watch.Event, strategy RouteStrategy) []string

type Route struct {
	sync.RWMutex
	routes map[string][]string
}
type Router struct {
	Route RouteFunc
}
