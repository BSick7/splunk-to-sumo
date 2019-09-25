package agent

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"io/ioutil"
	"strconv"
	"bytes"
)

func NewReverseProxy(sumoHost string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Transport: NewResponseCloseTransport(NewErrorTransport(http.DefaultTransport)),
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

			// we need to modify body to make it compatible with buffering
			body, _ := ioutil.ReadAll(req.Body)
			newbody := bytes.Replace(body,[]byte("}{"),[]byte("}\n{"),-1)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(newbody))
			req.ContentLength = int64(len(newbody))

			req.Header.Set("Content-Length", strconv.Itoa(len(newbody)))

			req.Header.Del("Authorization")
			req.Header.Del("Content-Type")
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		},
	}
}
