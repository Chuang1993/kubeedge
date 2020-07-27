package server

import (
	"net/http"
	"strings"

	"k8s.io/klog"

	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/util"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/server"
)

func WithRequestInfo(handler http.Handler) http.Handler {
	cfg := &server.Config{
		LegacyAPIGroupPrefixes: sets.NewString(server.DefaultLegacyAPIPrefix),
	}
	resolver := server.NewRequestInfoResolver(cfg)
	return filters.WithRequestInfo(handler, resolver)
}

func WithReqContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		reqContentType := request.Header.Get("Accept")
		parts := strings.Split(reqContentType, ",")
		if len(parts) >= 1 {
			reqContentType = parts[0]
		}
		if len(reqContentType) == 0 {
			klog.Errorf("request %s accept content type is null!", request.URL.String())
			http.Error(writer, "accept must need set.", http.StatusBadRequest)
			return
		}
		ctx = util.WithReqContentType(ctx, reqContentType)
		request = request.WithContext(ctx)
		handler.ServeHTTP(writer, request)
	})
}

func WithAppUserAgent(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userAgent := strings.ToLower(request.Header.Get("User-Agent"))
		parts := strings.Split(userAgent, "/")
		ua := "default"
		if len(parts) > 0 {
			ua = strings.ToLower(parts[0])
		}
		ctx = util.WithAppUserAgent(ctx, ua)
		request = request.WithContext(ctx)
		handler.ServeHTTP(writer, request)
	})
}