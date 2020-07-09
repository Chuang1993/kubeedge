package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/server"

	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/config"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/client"
)

func NewProxyServer() *ProxyServer {
	ps := &ProxyServer{
		mux: mux.NewRouter(),
		handler: &k8sHandler{
			metaClient: client.New(),
		},
	}
	return ps
}

type ProxyServer struct {
	mux     *mux.Router
	handler http.Handler
}

func (ps *ProxyServer) Run() {
	ps.installPath()

	server := &http.Server{
		Handler: ps.mux,
		Addr:    fmt.Sprintf(":%d", config.Config.ListenPort),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
func (ps *ProxyServer) installPath() {
	cfg := &server.Config{
		LegacyAPIGroupPrefixes: sets.NewString(server.DefaultLegacyAPIPrefix),
	}
	resolver := server.NewRequestInfoResolver(cfg)
	h := filters.WithRequestInfo(ps.handler, resolver)
	ps.mux.HandleFunc("/healthz", ps.healthz).Methods("GET")
	ps.mux.PathPrefix("/").Handler(h)
}

func (ps *ProxyServer) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
