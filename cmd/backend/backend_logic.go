package backend

import (
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	url           *url.URL
	alive         bool // to check if it is alive or dead
	current_conns int64 // keep tracking how many requests this backend handle
	mux           sync.RWMutex // mutex implementation to lock and lock for safety
}
func (b Backend) geturl() string{
	return b.url.String()
}
// healtch checker
func (b *Backend) setalive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	b.mux.Unlock()
}
func (b *Backend) isalive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive
}

func (b *Backend) increment_counter(){
	atomic.AddInt64(&b.current_conns,1)
}


func (b *Backend) decrement_counter(){
	atomic.AddInt64(&b.current_conns,-1)
}
