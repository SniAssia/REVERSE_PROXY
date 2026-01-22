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
func Newbackend(url1 string) (*Backend, error) {
    parsed, err := url.Parse(url1)
    if err != nil {
        return nil, err
    }
    b := &Backend{
        url: parsed,
    }
    b.alive= true
    return b, nil
}
func (b Backend) Activeconns() int64 {
	return b.current_conns
}
func (b Backend) Geturl() string{
	return b.url.String()
}
// healtch checker
func (b *Backend) Setalive(alive bool) {
	b.mux.Lock()
	b.alive = alive
	b.mux.Unlock()
}
func (b *Backend) Isalive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive
}

func (b *Backend) Increment_counter(){
	atomic.AddInt64(&b.current_conns,1)
}


func (b *Backend) Decrement_counter(){
	atomic.AddInt64(&b.current_conns,-1)
}
