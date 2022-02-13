package pertr

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

var ErrUsed = errors.New("context already used")

// options contains settings for wrapper.
type options struct {
	tr  http.RoundTripper
	ctx context.Context
}

type option func(*options)

// WithTransport sets upstream transport.
func WithTransport(tr http.RoundTripper) option {
	return func(o *options) {
		o.tr = tr
	}
}

// WithContext sets context to read body.
func WithContext(ctx context.Context) option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// wrapper represents http.RoundTrip.
type wrapper struct {
	sync.Mutex
	opts options
	used bool
}

// New creates new wrapper.
func New(opts ...option) http.RoundTripper {
	w := wrapper{
		opts: options{
			tr: http.DefaultTransport,
		},
	}
	for _, o := range opts {
		o(&w.opts)
	}
	return &w
}

func (w *wrapper) SetContext(ctx context.Context) {
	w.Lock()
	w.opts.ctx = ctx
	w.Unlock()
}

// RoundTrip executes http transaction and resume download operation if
// connection will be lost. It use Range header to start new download operations.
// RoundTrip() can't be used second time if context is used. If you want to use
// use transport second time, set new context with SetContext() first.
func (w *wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	rsp, err := w.opts.tr.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if rsp.Body != nil && rsp.Body != http.NoBody {
		w.Lock()
		defer w.Unlock()
		if w.used {
			return nil, ErrUsed
		}
		if w.opts.ctx != nil {
			w.used = true
		}
		rsp.Body = &body{rsp: rsp, up: rsp.Body, ctx: req.Context()}
	}

	return rsp, nil
}
