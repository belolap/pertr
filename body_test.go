package pertr

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testRsp struct {
	code   int
	header http.Header
	body   string
	close  bool
}

type testHandler struct {
	sync.Mutex
	server    *httptest.Server
	responses []testRsp
	iter      int
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Lock()
	defer h.Unlock()

	if h.iter >= len(h.responses) {
		panic("no more responses")
	}

	rsp := h.responses[h.iter]

	if rsp.code > 0 {
		w.WriteHeader(rsp.code)
	}

	if rsp.body != "" {
		_, err := w.Write([]byte(rsp.body))
		if err != nil {
			panic(err)
		}
	}

	if rsp.close {
		h.server.CloseClientConnections()
	}
}

type testChunk struct {
	num   int
	value string
}

func TestBody(t *testing.T) {
	tests := []struct {
		name string

		// Server settings.
		cfg       *http.Server // Server configuration.
		responses []testRsp    // Set of seq responses.

		// Request settings.
		method      string        // Request method.
		bodyTimeout time.Duration // Timeout to body read context.
		sizes       []int         // Body read buf sizes.

		// Expectations.
		expect    []testChunk // expect .Read() returns
		expectErr string
	}{
		{
			name:        "simple request",
			responses:   []testRsp{{code: http.StatusOK, body: "abc"}},
			method:      "GET",
			bodyTimeout: 1 * time.Second,
			sizes:       []int{3},
			expect:      []testChunk{{num: 3, value: "abc"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start test server with custom configuration.
			handler := testHandler{responses: tt.responses}
			ts := httptest.NewUnstartedServer(&handler)
			handler.server = ts
			if tt.cfg != nil {
				ts.Config = tt.cfg
			}
			ts.Start()
			defer ts.Close()

			// Set request context with download timeout.
			ctx := context.Background()
			if tt.bodyTimeout != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.bodyTimeout)
				defer cancel()
			}

			// Create client.
			cli := &http.Client{}
			cli.Transport = New(WithClient(cli), WithContext(ctx))

			// Create request.
			req, err := http.NewRequest(tt.method, ts.URL, http.NoBody)
			if err != nil {
				panic(err)
			}

			// Fire!
			rsp, err := cli.Do(req)
			if err != nil {
				panic(err)
			}
			defer rsp.Body.Close()

			// Parse results.
			err = nil
			for i, chunk := range tt.expect {
				buf := make([]byte, tt.sizes[i])
				l, err := rsp.Body.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						if chunk.num == l && i == len(tt.expect)-1 {
							// EOF at last chunk and all bytes recieved. No error.
							err = nil
						}
					}
					break
				}
				assert.Equal(t, chunk.num, l)
				assert.Equal(t, chunk.value, string(buf))
			}
			if tt.expectErr == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectErr)
				}
			}

		})
	}
}
