package backend

import (
	"net/url"
	"sync"
)

//managing all backends together

type Server struct {
	backends []*Backend
	current  uint64

	mux sync.RWMutex // always adding a mutex for the lock
}
func (s *Server) totalbackends() int {
	return len(s.backends)
}
func (s *Server) activebackends() int {
	return len(s.getalivebackends())
}


func (s *Server) addbackend(b *Backend) {
	s.mux.Lock()
	s.backends = append(s.backends, b)
	s.mux.Unlock()
}
func (s *Server) removebackend(url *url.URL) {
	s.mux.Lock()
	newBackends := s.backends[:0] 
	for _, b := range s.backends {
		if !(b.geturl()==url.String()) { 
			newBackends = append(newBackends, b)
		}
	}
s.backends = newBackends

	s.mux.Unlock()
}
func (s *Server) markbackendstatus(url *url.URL, alive bool){
	s.mux.Lock()
	defer s.mux.Unlock()
	for i := range s.backends{
		if s.backends[i].url == url {
			s.backends[i].setalive(alive)
		}
	}


}
func (s *Server) getalivebackends() []*Backend {
	var result []*Backend
	for i := range s.backends{
		if s.backends[i].isalive(){
			result = append(result , s.backends[i])
		}
	}
	return result
}
