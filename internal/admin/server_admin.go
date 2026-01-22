package admin

import (
	"REVERSE_PROXY/cmd/backend"
	"net/http"
	"sync/atomic"
	"context"
)

// this file will basically do : Creating the HTTP server
// Wiring routes to handlers
// Starting and stopping the admin server
type Admin struct {
	server       *http.Server
	pool         *backend.Server
	loadbalancer *atomic.Value // holds LoadBalancer
}

func (a *Admin) Start() error {
	return a.server.ListenAndServe()
}
func newserver(address string, pool *backend.Server, loadbal *atomic.Value) *Admin{
	s := &Admin {
		pool : pool,
		loadbalancer: loadbal,

	}
	mux := http.NewServeMux()
	mux.HandleFunc("/status",s.handleStatus)
	mux.HandleFunc("/backends",s.handleBackends)
	mux.HandleFunc("/strategy",s.handleStrategy)
	s.server = &http.Server {
		Addr : address,
		Handler : mux ,
	}
	return s 
}

func (s *Admin) Shutdown(ctx context.Context) error {
    return s.server.Shutdown(ctx)
}