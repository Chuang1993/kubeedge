package registry

import (
	"context"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"

	edgewatch "github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/watch"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/client"
)

func init() {
	sms := &serviceMetaStore{
		metaClient: client.New(),
	}
	RegistryGetter("services", sms)
	RegistryLister("services", sms)
	RegistryWatcher("services", sms)
}

type serviceMetaStore struct {
	metaClient client.CoreInterface
}

func (sms *serviceMetaStore) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	reqInfo, _ := apirequest.RequestInfoFrom(ctx)
	svc, err := sms.metaClient.Services(reqInfo.Namespace).Get(name)
	return svc, err
}

func (sms *serviceMetaStore) NewList() runtime.Object {
	return &v1.ServiceList{}
}

func (sms *serviceMetaStore) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	reqInfo, _ := apirequest.RequestInfoFrom(ctx)
	objs, err := sms.metaClient.Services("").ListAll()
	if err != nil {
		return nil, err
	}
	listRv := 0
	rvStr := ""
	rvInt := 0
	accessor := meta.NewAccessor()
	for i := range objs {
		rvStr, _ = accessor.ResourceVersion(&objs[i])
		rvInt, _ = strconv.Atoi(rvStr)
		if rvInt > listRv {
			listRv = rvInt
		}
	}

	svcList := &v1.ServiceList{
		TypeMeta: metav1.TypeMeta{},
		ListMeta: metav1.ListMeta{ResourceVersion: strconv.Itoa(listRv)},
		Items:    objs,
	}
	accessor.SetKind(svcList, "ServiceList")
	accessor.SetAPIVersion(svcList, reqInfo.APIVersion)
	return svcList, nil
}

func (sms *serviceMetaStore) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	watcher := edgewatch.Incoming(nil)
	rv, _ := strconv.ParseUint(options.ResourceVersion, 0, 64)
	//edgewatch.AddWatcher("service", uuid.NewV4().String(), rv, watcher, time.Duration(*options.TimeoutSeconds)*time.Second)
	edgewatch.AddWatcher("service", uuid.NewV4().String(), rv, watcher, 400*time.Second)
	// TODO 构建过滤条件
	filterWatcher := watch.Filter(watcher, edgewatch.AlwaysTrue)
	return filterWatcher, nil
}
