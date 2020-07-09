package registry

import (
	"context"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

//func init() {
//	pms := &podMetaStore{}
//	RegistryGetter("pod", pms)
//	RegistryLister("pod", pms)
//	RegistryWatcher("pod", pms)
//}

type podMetaStore struct {
}

func (pms *podMetaStore) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, nil
}

func (pms *podMetaStore) NewList() runtime.Object {
	return nil
}

func (pms *podMetaStore) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	return nil, nil
}

func (pms *podMetaStore) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	return nil, nil
}
