package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kubeedge/kubeedge/edge/pkg/edgeproxy/config"
	"github.com/kubeedge/kubeedge/edge/pkg/metamanager/client"
	"net/http"
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
	ps.mux.HandleFunc("/healthz", ps.healthz).Methods("GET")
	ps.mux.PathPrefix("/").Handler(ps.handler)
}

func (ps *ProxyServer) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
