package backend

import (
	"net/url"
	"sync"
)

//managing all backends together

type server struct {
	backends []*Backend
	current  uint64

	mux sync.RWMutex // always adding a mutex for the lock
}
func (s *server) totalbackends() int {
	return len(s.backends)
}
func (s *server) activebackends() int {
	return len(s.getalivebackends())
}


func (s *server) addbackend(b *Backend) {
	s.mux.Lock()
	s.backends = append(s.backends, b)
	s.mux.Unlock()
}
func (s *server) removebackend(url *url.URL) {
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
func (s *server) markbackendstatus(url *url.URL, alive bool){
	s.mux.Lock()
	for i := range s.backends{
		if s.backends[i].url == url {
			s.backends[i].setalive(alive)
		}
	}

}
func (s *server) getalivebackends() []*Backend {
	result := []*Backend{}
	for i := range s.backends{
		if s.backends[i].isalive() == true {
			result = append(result , s.backends[i])
		}
	}
	return result
}
