package pertr

import (
	"net/http"
)

// Wrapper represents http.RoundTrip.
type Wrapper struct {
	Transport http.RoundTripper
}

// RoundTrip executes http transaction and resume download operation if
// connection will be lost. It use Range header to start new download operations.
func (w Wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	var t http.RoundTripper
	if w.Transport != nil {
		t = w.Transport
	} else {
		t = http.DefaultTransport
	}

	rsp, err := t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	rsp.Body = &body{rsp: rsp, up: rsp.Body, ctx: req.Context()}

	return rsp, nil
}
