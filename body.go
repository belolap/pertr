package pertr

import (
	"context"
	"io"
	"net/http"
)

// body reads http response and reconnects if connection has been lost.
type body struct {
	rsp *http.Response
	up  io.ReadCloser
	ctx context.Context
}

// Read http response stream.
func (b *body) Read(p []byte) (n int, err error) {
	return b.up.Read(p)
}

// Close http response stream.
func (b *body) Close() error {
	return b.up.Close()
}
