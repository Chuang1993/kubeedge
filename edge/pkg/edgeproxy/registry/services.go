package registry

import (
	"context"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type serviceMetaStore struct {
}

func (pms *serviceMetaStore) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, nil
}

func (pms *serviceMetaStore) NewList() runtime.Object {
	return nil
}

func (pms *serviceMetaStore) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	return nil, nil
}

func (pms *serviceMetaStore) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	return nil, nil
}
