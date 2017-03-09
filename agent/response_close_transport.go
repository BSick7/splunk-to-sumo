package agent

import "net/http"

type ResponseCloseTransport struct {
	Transport http.RoundTripper
}

func NewResponseCloseTransport(transport http.RoundTripper) *ResponseCloseTransport {
	return &ResponseCloseTransport{Transport: transport}
}

func (t *ResponseCloseTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(request)
	if err != nil {
		return resp, err
	}
	return resp, resp.Body.Close()
}
