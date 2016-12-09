package agent

import (
	"net/http"
	"net/http/httputil"
)

type Server struct {
	rp *httputil.ReverseProxy
}

func NewServer(sumoHost string) *Server {
	return &Server{
		rp: NewReverseProxy(sumoHost),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		// used by splunk logging driver to verify connection
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "POST" && r.URL.Path == "/services/collector/event/1.0" {
		s.rp.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
