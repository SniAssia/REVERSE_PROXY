package admin
import ("net/http"
		"encoding/json"
    "REVERSE_PROXY/internal/loadbalancer"
    "REVERSE_PROXY/cmd/backend")
		

func (s *Admin) handleStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    lb := s.loadbalancer.Load().(loadbalancer.LoadBalancer)
    

    backends := s.pool.Backends()
    var status []Backendstatus

    for _, b := range backends {
        status = append(status, Backendstatus{
            URL:   b.Geturl(),
            Alive: b.Isalive(),
            Conns: b.Activeconns(), 
        })
    }
    resp := Statusresponse{
        Strategy: lb.Name(),
        Backends: status,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (s *Admin) handleBackends(w http.ResponseWriter, r *http.Request) {
    var req Backendrequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    switch r.Method {

    case http.MethodPost:
        backend, err := backend.Newbackend(req.URL)
        if err != nil {
            http.Error(w, "invalid backend url", http.StatusBadRequest)
            return
        }

        s.pool.Addbackend(backend)
        w.WriteHeader(http.StatusCreated)

    case http.MethodDelete:
        removed := s.pool.Removebackend(req.URL)
        if !removed {
            http.Error(w, "backend not found", http.StatusNotFound)
            return
        }

        w.WriteHeader(http.StatusOK)

    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func (s *Admin) handleStrategy(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req Strategyrequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    var lb loadbalancer.LoadBalancer

    switch req.Strategy {
    case "round-robin":
        lb = loadbalancer.Newround(s.pool)
    case "least-connections":
        lb = loadbalancer.Newleastconn(s.pool)
    default:
        http.Error(w, "unknown strategy", http.StatusBadRequest)
        return
    }

    s.loadbalancer.Store(lb)
    w.WriteHeader(http.StatusOK)
}

