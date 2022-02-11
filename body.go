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

// Read http response stream. If connection will lost reading will be resumed with
// new http request contains Range header. Connection lost detects by folowing signs:
//   - upstream returns EOF and request number of bytes != recieved number of bytes
//   - upsteram returns context timeout exceeded error
func (b *body) Read(p []byte) (n int, err error) {
	// Handle errors:
	// - unexpected EOF
	// - context deadline exeeded
	//
	// detect errors
	return b.up.Read(p)
}

// Close http response stream.
func (b *body) Close() error {
	return b.up.Close()
}
