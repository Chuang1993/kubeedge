package server

import (
	"net/http"
	"path"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers"
	"k8s.io/apiserver/pkg/endpoints/request"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/registry"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/client"
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
	ctx := r.Context()
	reqinfo, _ := apirequest.RequestInfoFrom(ctx)
	scope := handler.getRequestScope(reqinfo)
	lister := registry.GetLister(reqinfo.Resource)
	watcher := registry.GetWatcher(reqinfo.Resource)
	handlers.ListResource(lister, watcher, scope, true, 5*time.Minute).ServeHTTP(w, r)
}

func (handler *k8sHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqinfo, _ := apirequest.RequestInfoFrom(ctx)
	scope := handler.getRequestScope(reqinfo)
	lister := registry.GetLister(reqinfo.Resource)
	handlers.ListResource(lister, nil, scope, false, 0).ServeHTTP(w, r)
}

func (handler *k8sHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqinfo, _ := apirequest.RequestInfoFrom(ctx)
	scope := handler.getRequestScope(reqinfo)
	getter := registry.GetGetter(reqinfo.Resource)
	handlers.GetResource(getter, nil, scope).ServeHTTP(w, r)
}

func (handler *k8sHandler) getRequestScope(reqinfo *request.RequestInfo) *handlers.RequestScope {
	clusterScoped := reqinfo.Namespace == ""
	prefix := "/" + path.Join(reqinfo.APIPrefix, reqinfo.APIGroup, reqinfo.APIVersion)
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
	return scope
}
