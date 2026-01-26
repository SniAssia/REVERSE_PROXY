package main

import (
	
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"REVERSE_PROXY/cmd/backend"
	"REVERSE_PROXY/internal/loadbalancer"
)

// Proxyhandler routes requests to backends using a Loadbalancer
type ProxyHandler struct {
	server *backend.Server
	lb     *atomic.Value // holds LoadBalancer
}

// Constructor
func NewProxyHandler(server *backend.Server, lb *atomic.Value) *ProxyHandler {
	return &ProxyHandler{
		server: server,
		lb:     lb,
	}
}
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb := p.lb.Load().(loadbalancer.LoadBalancer)

	backend, err := lb.NextBackend(p.server)
	if err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	backend.Increment_counter()
	defer backend.Decrement_counter()
	targetStr := backend.Geturl()          
	targetURL, err := url.Parse(targetStr) 
	if err != nil{
		
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, e error) {
		log.Printf("[proxy] backend %s failed: %v", targetStr, e)
		backend.Setalive(false)
		http.Error(resp, "Service Unavailable", http.StatusServiceUnavailable)
	}

	proxy.ServeHTTP(w, r)
}
