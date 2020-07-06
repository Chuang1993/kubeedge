package util

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
)

func WriteObject(statusCode int, obj runtime.Object, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gv := schema.GroupVersion{
		Group:   "",
		Version: runtime.APIVersionInternal,
	}
	if info, ok := apirequest.RequestInfoFrom(ctx); ok {
		gv.Group = info.APIGroup
		gv.Version = info.APIVersion
	}

	responsewriters.WriteObjectNegotiated(clientscheme.Codecs, negotiation.DefaultEndpointRestrictions, gv, w, r, statusCode, obj)
}

// Err write err to response writer
func Err(err error, w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	gv := schema.GroupVersion{
		Group:   "",
		Version: runtime.APIVersionInternal,
	}
	if info, ok := apirequest.RequestInfoFrom(ctx); ok {
		gv.Group = info.APIGroup
		gv.Version = info.APIVersion
	}

	responsewriters.ErrorNegotiated(err, clientscheme.Codecs, gv, w, req)
}
