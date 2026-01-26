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
// in order to avoid deadlocks and race conditions 
func (s *Server) Backends() []*Backend {
    s.mux.RLock()
    defer s.mux.RUnlock()

    snapshot := make([]*Backend, len(s.backends))
    copy(snapshot, s.backends)
    return snapshot
}
func (s *Server) Totalbackends() int {
	return len(s.backends)
}
func (s *Server) Activebackends() int {
	return len(s.Getalivebackends())
}


func (s *Server) Addbackend(b *Backend) {
	s.mux.Lock()
	s.backends = append(s.backends, b)
	s.mux.Unlock()
}
func (s *Server) Removebackend(url string) bool {
	s.mux.Lock()
	removed := &Backend{}
	newBackends := s.backends[:0] 
	for _, b := range s.backends {
		if !(b.Geturl()==url) { 
			removed = b 
			newBackends = append(newBackends, b)
		}
		removed = nil
	}
	s.backends = newBackends

	s.mux.Unlock()
	if removed == nil {
		return false 
	} else {
		return true 
	}
}
func (s *Server) Markbackendstatus(url *url.URL, alive bool){
	s.mux.Lock()
	defer s.mux.Unlock()
	for i := range s.backends{
		if s.backends[i].url == url {
			s.backends[i].Setalive(alive)
		}
	}


}
func (s *Server) Getalivebackends() []*Backend {
	var result []*Backend
	for i := range s.backends{
		if s.backends[i].Isalive(){
			result = append(result , s.backends[i])
		}
	}
	return result
}
func NewServer() *Server {
    return &Server{
    }
}