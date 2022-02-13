package pertr

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// body reads http response and reconnects if connection has been lost.
type body struct {
	sync.Mutex

	tr  http.RoundTripper
	rsp *http.Response
	up  io.ReadCloser
	ctx context.Context

	readed int64
	ferr   error // Save first error.
}

// Read http response stream. If connection will lost reading will be resumed with
// new http request contains Range header.
func (b *body) Read(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()

	n, err = b.up.Read(p)
	b.readed += int64(n)

	if err == nil || err == io.EOF {
		return
	}

	if b.ferr == nil {
		b.ferr = err
	}

	_ = b.up.Close()

	select {
	case <-b.ctx.Done():
		return n, context.DeadlineExceeded
	default:
	}

	req := b.rsp.Request.Clone(b.rsp.Request.Context())
	req.Header.Set("Range", fmt.Sprintf("%d-", b.readed))

	rsp, err := b.tr.RoundTrip(req)
	if err != nil {
		return n, err
	}

	// Server must set Content-Range header, just not care about it.
	// Server must set Content-Length with full body's length - ignore.
	// Server must set Content-Type to multipart/byteranges - just not care.
	// Server must return http.StatusPartial content.
	if rsp.StatusCode == http.StatusPartialContent {
		if rsp.Body == nil || rsp.Body == http.NoBody {
			// Server not returns body, just return first error.
			return n, b.ferr
		}
		b.up = rsp.Body
		return n, nil
	}

	// Server returns either http.StatusOK that means server not supports resume
	// or another status so we don't know what to do. So just return first error.
	return n, b.ferr
}

// Close http response stream.
func (b *body) Close() error {
	return b.up.Close()
}
