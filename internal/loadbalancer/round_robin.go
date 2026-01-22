package loadbalancer

import ("REVERSE_PROXY/cmd/backend"
		"sync/atomic")
		

type RoundRobin struct {
	pool  *backend.Server
	current uint64
}

func (r *RoundRobin) Name() string {
	return "round-robin"
}
func Newround(pool *backend.Server) *RoundRobin{
	var1 := &RoundRobin{
		pool : pool,
	}
	return var1

}
func (r *RoundRobin) NextBackend(pool *backend.Server) (*backend.Backend, error) {
	backends := r.pool.Backends()
	n := r.pool.Totalbackends()
	if n== 0 {
		return nil,ErrNoAvailableBackends
	}
	// here i use atomic op so i make sure that this step is 
	// indivisible 
	start := atomic.AddUint64(&r.current, 1)

    for i := 0; i < n; i++ {
        idx := int((start + uint64(i)) % uint64(n))
        backend := backends[idx]

        if backend.Isalive() {
            return backend, nil
        }
    }

    return nil, ErrNoAvailableBackends


}