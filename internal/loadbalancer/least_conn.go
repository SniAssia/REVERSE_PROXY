package loadbalancer

import ("math"
"REVERSE_PROXY/cmd/backend") 

type least_conn struct {
	pool *backend.Server 

}
func Newleastconn(pool *backend.Server) *least_conn {
	var1 := least_conn{
		pool : pool ,
	}
	return &var1
}
func (l *least_conn) Name() string {
	return "least-connections"
}
func (l *least_conn) NextBackend(*backend.Server) (*backend.Backend,error) {
	backends := l.pool.Backends()
	var selected *backend.Backend
	min1 := int64(math.MaxInt64)

	for _,backend := range backends {
		if !backend.Isalive(){
			continue
		}
		conns := backend.Activeconns()
		if conns < min1 {
			min1 = conns 
			selected = backend
		}
	}
	if selected == nil {
		return nil , ErrNoAvailableBackends
	}
	return selected , nil 

}