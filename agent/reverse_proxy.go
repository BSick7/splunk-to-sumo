package agent

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewReverseProxy(sumoHost string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Transport: NewErrorTransport(http.DefaultTransport),
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
