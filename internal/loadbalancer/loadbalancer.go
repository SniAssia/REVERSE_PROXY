package loadbalancer

import "REVERSE_PROXY/cmd/backend"
type LoadBalancer interface {
    NextBackend(pool *backend.Server) (*backend.Backend, error)
    Name() string
}
