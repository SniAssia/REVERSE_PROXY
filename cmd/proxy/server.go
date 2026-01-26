package main

import (
	"context"
	"net/http"
	"sync/atomic"

	"REVERSE_PROXY/cmd/backend"
)

// ProxyServer wraps an HTTP server + backend pool + load balancer
type ProxyServer struct {
	server *http.Server
	pool   *backend.Server
	lb     *atomic.Value
}

// Constructor
func NewProxyServer(addr string, pool *backend.Server, lb *atomic.Value) *ProxyServer {
	handler := NewProxyHandler(pool, lb)

	return &ProxyServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		pool: pool,
		lb:   lb,
	}
}

// Start listening
func (p *ProxyServer) ListenAndServe() error {
	return p.server.ListenAndServe()
}

// Graceful shutdown
func (p *ProxyServer) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}
