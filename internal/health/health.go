package health

import (
	"context"
	"log"
	"net/http"
	"time"

	"REVERSE_PROXY/cmd/backend"
)

type Checker struct {
	pool     *backend.Server
	interval time.Duration
	client   *http.Client
}

func NewChecker(pool *backend.Server, interval time.Duration) *Checker {
	return &Checker{
		pool:     pool,
		interval: interval,
		client:   &http.Client{Timeout: 2 * time.Second},
	}
}

func (h *Checker) Start(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.checkAllBackends()
		case <-ctx.Done():
			log.Println("[health] stopping health checker")
			return
		}
	}
}

func (h *Checker) checkAllBackends() {
	backends := h.pool.Backends()
	for _, b := range backends {
		alive := h.pingBackend(b)
		if alive != b.Isalive() {
			b.Setalive(alive)
			log.Printf("[health] backend %s alive=%v\n", b.Geturl(), alive)
		}
	}
}

func (h *Checker) pingBackend(b *backend.Backend) bool {
	resp, err := h.client.Get(b.Geturl() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
