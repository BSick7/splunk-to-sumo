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
	if len(os.Args) != 3 {
		usage()
		os.Exit(0)
	}

	bind := os.Args[1]
	host := os.Args[2]

	server := newSplunkToSumoServer(host)
	fmt.Printf("Starting splunk-to-sumo v%s on %s with hostname=%s...\n", Version, bind, host)
	http.ListenAndServe(fmt.Sprintf(bind), server)
}

func usage() {
	fmt.Println(`usage: splunk-to-sumo [bind-address] [sumo-host]

  bind-address      This specifies which ip:port to bind the server.

  sumo-host         This will forward messages to sumologic with
                    this set as the X-Sumo-Host header.
`)
}

type splunkToSumoServer struct {
	rp *httputil.ReverseProxy
}

func newSplunkToSumoServer(sumoHost string) *splunkToSumoServer {
	return &splunkToSumoServer{
		rp: newSplunkToSumoReverseProxy(sumoHost),
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

func newSplunkToSumoReverseProxy(sumoHost string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Authorization header is packed with sumologic data
			auth := req.Header.Get("Authorization")
			rawUrl := strings.TrimPrefix(auth, "Splunk ")
			endpoint, _ := url.Parse(rawUrl)
			req.Header.Add("X-Sumo-Name", endpoint.Query().Get("name"))
			req.Header.Add("X-Sumo-Category", endpoint.Query().Get("category"))
			req.Header.Add("X-Sumo-Host", sumoHost)
			endpoint.RawQuery = ""
			req.URL = endpoint

			req.Header.Del("Authorization")
			req.Header.Del("Content-Type")
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		},
	}
}
