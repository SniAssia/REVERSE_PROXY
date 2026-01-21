package loadbalancer

//import "REVERSE_PROXY/cmd/backend"

type RoundRobin struct {
	current uint64
}

func (rr *RoundRobin) Name() string {
	return "round-robin"
}
/**
func (rr *RoundRobin) NextBackend(pool *backend.Server) (*backend.Backend,error){
	

}**/
