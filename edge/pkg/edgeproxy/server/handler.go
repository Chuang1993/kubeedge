package server

import (
	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/registry"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/client"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"path"
	"time"
)

var resourceToKind = map[string]string{
	"pod":       "PodList",
	"service":   "ServiceList",
	"endpoints": "EndpointsList",
}

type k8sHandler struct {
	metaClient client.CoreInterface
}

func (handler *k8sHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO need support subreource?
	if reqInfo, ok := apirequest.RequestInfoFrom(ctx); ok && reqInfo != nil && reqInfo.IsResourceRequest {
		switch reqInfo.Verb {
		case "watch":
			handler.watch(w, r)
		case "list":
			handler.list(w, r)
		case "get":
			handler.get(w, r)
		default:
			// update operation will return the newest resource object , we don't update the cloud side object ,so we just return the metaManager object。\
			// if the operation is create,we should return err
			handler.get(w, r)
		}
	}
}
func (handler *k8sHandler) watch(w http.ResponseWriter, r *http.Request) {
	// 构建watch的channel，channel 应统一管理，以便推送新的数据过来，
	//注意判断resouceVersion,应先查询一次，然后进入通道采集数据阶段，避免在list和watch请求间隔之间存在数据丢失。
	//或者在list阶段就构建这部分内容。watch阶段直接使用，定义好生存时间即可。 如果watch操作没有找到指定的内容，则可直接返回异常，客户端再次发起list/watch请求
	ctx := r.Context()
	reqinfo, _ := apirequest.RequestInfoFrom(ctx)
	// 构建lister，watcher接口，维护watcher接口
	clusterScoped := reqinfo.Namespace == ""
	prefix := "/" + path.Join(reqinfo.APIGroup, reqinfo.APIGroup)
	namer := handlers.ContextBasedNaming{
		SelfLinker:         runtime.SelfLinker(meta.NewAccessor()),
		SelfLinkPathPrefix: path.Join(prefix, reqinfo.Resource) + "/",
		SelfLinkPathSuffix: "",
		ClusterScoped:      clusterScoped,
	}

	scope := &handlers.RequestScope{
		Namer:               namer,
		Serializer:          scheme.Codecs,
		ParameterCodec:      scheme.ParameterCodec,
		StandardSerializers: scheme.Codecs.SupportedMediaTypes(),
		Creater:             scheme.Scheme,
		Convertor:           scheme.Scheme,
		Defaulter:           scheme.Scheme,
		Typer:               scheme.Scheme,
		UnsafeConvertor:     scheme.Scheme,
		Authorizer:          nil,
		Resource:            schema.GroupVersionResource{Group: reqinfo.APIGroup, Version: reqinfo.APIVersion, Resource: reqinfo.Resource},
		Kind:                schema.GroupVersionKind{Group: reqinfo.APIGroup, Version: reqinfo.APIVersion, Kind: reqinfo.Resource},
		Subresource:         "",
		MaxRequestBodyBytes: int64(3 * 1024 * 1024),
	}
	lister := registry.GetLister(reqinfo.Resource)
	watcher := registry.GetWatcher(reqinfo.Resource)
	handlers.ListResource(lister, watcher, scope, true, 5*time.Minute).ServeHTTP(w, r)
	//gv := &schema.GroupVersion{Group: reqinfo.APIGroup, Version: reqinfo.APIVersion}
	//serializer, err := negotiation.NegotiateOutputMediaTypeStream(r, scheme.Codecs, negotiation.DefaultEndpointRestrictions)
	//if err != nil {
	//	util.Err(err, w, r)
	//	return
	//}
	//framer := serializer.StreamSerializer.Framer
	//encoder := scheme.Codecs.EncoderForVersion(serializer.StreamSerializer.Serializer, gv)
	//flusher, ok := w.(http.Flusher)
	//if !ok {
	//	err := fmt.Errorf("unable to start watch - can't get http.Flusher: %#v", w)
	//	utilruntime.HandleError(err)
	//	util.Err(err, w, r)
	//	return
	//}
	//w.Header().Set("Content-Type", serializer.MediaType)
	//w.Header().Set("Transfer-Encoding", "chunked")
	//w.WriteHeader(http.StatusOK)
	//flusher.Flush()
	//fw := framer.NewFrameWriter(w)
	//e := streaming.NewEncoder(fw, encoder)
	//internalEvent := &metav1.InternalEvent{}
	//outEvent := &metav1.WatchEvent{}
	//buf := &bytes.Buffer{}
	//ch := make(chan watch.Event)
	//done := r.Context().Done()
	//for {
	//	//
	//	select {
	//	case <-done:
	//		return
	//	case event, ok := <-ch:
	//		if !ok {
	//			return
	//		}
	//		fmt.Sprintf("%v", event)
	//	}
	//}
	//return
}

func (handler *k8sHandler) list(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	reqinfo, _ := apirequest.RequestInfoFrom(ctx)
	clusterScoped := reqinfo.Namespace == ""
	prefix := "/" + path.Join(reqinfo.APIGroup, reqinfo.APIGroup)
	namer := handlers.ContextBasedNaming{
		SelfLinker:         runtime.SelfLinker(meta.NewAccessor()),
		SelfLinkPathPrefix: path.Join(prefix, reqinfo.Resource) + "/",
		SelfLinkPathSuffix: "",
		ClusterScoped:      clusterScoped,
	}

	scope := &handlers.RequestScope{
		Namer:               namer,
		Serializer:          scheme.Codecs,
		ParameterCodec:      scheme.ParameterCodec,
		StandardSerializers: scheme.Codecs.SupportedMediaTypes(),
		Creater:             scheme.Scheme,
		Convertor:           scheme.Scheme,
		Defaulter:           scheme.Scheme,
		Typer:               scheme.Scheme,
		UnsafeConvertor:     scheme.Scheme,
		Authorizer:          nil,
		Resource:            schema.GroupVersionResource{Group: reqinfo.APIGroup, Version: reqinfo.APIVersion, Resource: reqinfo.Resource},
		Kind:                schema.GroupVersionKind{Group: reqinfo.APIGroup, Version: reqinfo.APIVersion, Kind: reqinfo.Resource},
		Subresource:         "",
		MaxRequestBodyBytes: int64(3 * 1024 * 1024),
	}
	lister := registry.GetLister(reqinfo.Resource)
	handlers.ListResource(lister, nil, scope, false, 0).ServeHTTP(w, r)
	// TODO 构造list结构体，返回最新的resouceversion即可，构建selfLink
	//ctx := r.Context()
	//reqinfo, _ := apirequest.RequestInfoFrom(ctx)
	//listKind := resourceToKind[reqinfo.Resource]
	//gkv := schema.GroupVersionKind{
	//	Group:   reqinfo.APIGroup,
	//	Version: reqinfo.APIVersion,
	//	Kind:    listKind,
	//}
	//opts := metainternalversion.ListOptions{}
	//values := r.URL.Query()
	//if len(values) > 0 {
	//	metainternalversionscheme.ParameterCodec.DecodeParameters(values, schema.GroupVersion{Version: reqinfo.APIVersion, Group: reqinfo.APIGroup}, &opts)
	//}
	////svcs, err := handler.metaClient.Services("all").ListAll()
	////if err != nil {
	////	util.Err(err, w, r)
	////}
	//// TODO 从metaManager内查询逻辑
	//objs := make([]runtime.Object, 0)
	//listRv := 0
	//rvStr := ""
	//rvInt := 0
	//accessor := meta.NewAccessor()
	//for i := range objs {
	//	rvStr, _ = accessor.ResourceVersion(objs[i])
	//	rvInt, _ = strconv.Atoi(rvStr)
	//	if rvInt > listRv {
	//		listRv = rvInt
	//	}
	//}
	//listobj, err := scheme.Scheme.New(gkv)
	//if err != nil {
	//	util.Err(err, w, r)
	//	return
	//}
	//accessor.SetResourceVersion(listobj, strconv.Itoa(listRv))
	//clusterScoped := true
	//if reqinfo.Namespace != "" {
	//	clusterScoped = false
	//}
	//
	//prefix := "/" + path.Join(reqinfo.APIGroup, reqinfo.APIGroup)
	//namer := handlers.ContextBasedNaming{
	//	SelfLinker:         runtime.SelfLinker(meta.NewAccessor()),
	//	SelfLinkPathPrefix: path.Join(prefix, reqinfo.Resource) + "/",
	//	SelfLinkPathSuffix: "",
	//	ClusterScoped:      clusterScoped,
	//}
	//
	//uri, err := namer.GenerateListLink(r)
	//if err != nil {
	//	util.Err(err, w, r)
	//	return
	//}
	//if err := namer.SetSelfLink(listobj, uri); err != nil {
	//	util.Err(err, w, r)
	//	return
	//}
	//util.WriteObject(http.StatusOK, listobj, w, r)
}

func (handler *k8sHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqinfo, _ := apirequest.RequestInfoFrom(ctx)
	clusterScoped := reqinfo.Namespace == ""
	prefix := "/" + path.Join(reqinfo.APIGroup, reqinfo.APIGroup)
	namer := handlers.ContextBasedNaming{
		SelfLinker:         runtime.SelfLinker(meta.NewAccessor()),
		SelfLinkPathPrefix: path.Join(prefix, reqinfo.Resource) + "/",
		SelfLinkPathSuffix: "",
		ClusterScoped:      clusterScoped,
	}

	scope := &handlers.RequestScope{
		Namer:               namer,
		Serializer:          scheme.Codecs,
		ParameterCodec:      scheme.ParameterCodec,
		StandardSerializers: scheme.Codecs.SupportedMediaTypes(),
		Creater:             scheme.Scheme,
		Convertor:           scheme.Scheme,
		Defaulter:           scheme.Scheme,
		Typer:               scheme.Scheme,
		UnsafeConvertor:     scheme.Scheme,
		Authorizer:          nil,
		Resource:            schema.GroupVersionResource{Group: reqinfo.APIGroup, Version: reqinfo.APIVersion, Resource: reqinfo.Resource},
		Kind:                schema.GroupVersionKind{Group: reqinfo.APIGroup, Version: reqinfo.APIVersion, Kind: reqinfo.Resource},
		Subresource:         "",
		MaxRequestBodyBytes: int64(3 * 1024 * 1024),
	}
	getter := registry.GetGetter(reqinfo.Resource)
	handlers.GetResource(getter, nil, scope).ServeHTTP(w, r)
	//var obj runtime.Object
	//switch reqinfo.Resource {
	//case "Pod":
	//	obj, err := handler.metaClient.Pods(reqinfo.Namespace).Get(reqinfo.Name)
	//	if err != nil {
	//		return err
	//	}
	//case "Service":
	//	obj, err := handler.metaClient.Services(reqinfo.Namespace).Get(reqinfo.Name)
	//	if err != nil {
	//		return err
	//	}
	//}
	// TODO 从metaManager内无差异的获取数据，可以传递namespace,name,resouce
	//obj, err := handler.metaClient.Services(reqinfo.Namespace).Get(reqinfo.Name)
	//if err != nil {
	//	util.Err(err, w, r)
	//	return
	//}
	//util.WriteObject(http.StatusOK, obj, w, r)
}
