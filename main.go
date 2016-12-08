package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var Version string

func main() {
	bind := "0.0.0.0:5001"
	if len(os.Args) > 0 {
		bind = os.Args[0]
	}

	server := newSplunkToSumoServer()
	fmt.Println("Starting splunk-to-sumo %s on %s", Version, bind)
	http.ListenAndServe(fmt.Sprintf(bind), server)
}

type splunkToSumoServer struct {
	rp *httputil.ReverseProxy
}

func newSplunkToSumoServer() *splunkToSumoServer {
	return &splunkToSumoServer{
		rp: newSplunkToSumoReverseProxy(),
	}
}

func (s *splunkToSumoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "POST" && r.URL.Path == "/services/collector/event/1.0" {
		s.rp.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func newSplunkToSumoReverseProxy() *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Authorization header is packed with sumologic data
			auth := req.Header.Get("Authorization")
			rawUrl := strings.TrimPrefix(auth, "Splunk ")
			endpoint, _ := url.Parse(rawUrl)
			req.Header.Add("X-Sumo-Category", endpoint.Query().Get("category"))
			req.Header.Add("X-Sumo-Name", endpoint.Query().Get("name"))
			endpoint.RawQuery = ""
			req.URL = endpoint

			fmt.Println(req.Header)

			fmt.Printf("Forwarding to %s\n", req.URL)

			req.Header.Del("Authorization")
			req.Header.Del("Content-Type")
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		},
	}
}
