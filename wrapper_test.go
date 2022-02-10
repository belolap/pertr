package pertr

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testModeHeader = "X-TestRequestMode"

type testResponse struct{}

func (rsp *testResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mode := r.Header.Get(testModeHeader)
	switch mode {

	case "echo":
		// Simple echo.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			panic(err)
		}

	default:
		// Unknown mode.
		panic("unknown test mode")
	}
}

// testChunk contains read response body chunk settings.
type testChunk struct {
	c int    // cap, size of buffer to read response body chunk
	l int    // len, number of real returned bytes
	v string // returned value
}

func TestTransport(t *testing.T) {
	tests := []struct {
		name      string
		mode      string
		cfg       *http.Server
		method    string
		body      string
		timeout   time.Duration
		expect    []testChunk
		expecterr string
	}{
		{
			name:    "simple request",
			mode:    "echo",
			method:  "get",
			timeout: 100 * time.Second,
			body:    "abc",
			expect:  []testChunk{{c: 3, l: 3, v: "abc"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start test server with custom configuration.
			ts := httptest.NewUnstartedServer(&testResponse{})
			if tt.cfg != nil {
				ts.Config = tt.cfg
			}
			ts.Start()
			defer ts.Close()

			// Set request context with download timeout.
			ctx := context.Background()
			if tt.timeout != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.timeout)
				defer cancel()
			}

			// Set request body.
			var body io.Reader
			if tt.body != "" {
				body = bytes.NewReader([]byte(tt.body))
			} else {
				body = http.NoBody
			}

			// Create client.
			cli := http.Client{Transport: Wrapper{}}

			// Create request.
			req, err := http.NewRequestWithContext(ctx, tt.method, ts.URL, body)
			if err != nil {
				panic(err)
			}
			req.Header.Add(testModeHeader, tt.mode)

			// Fire!
			rsp, err := cli.Do(req)
			if err != nil {
				panic(err)
			}
			defer rsp.Body.Close()

			// Parse results.
			for i, chunk := range tt.expect {
				buf := make([]byte, chunk.c)
				l, err := rsp.Body.Read(buf)
				if tt.expecterr == "" {
					if errors.Is(err, io.EOF) {
						if chunk.c == l && i == len(tt.expect)-1 {
							continue
						}
					} else if assert.NoError(t, err) {
						assert.Equal(t, chunk.l, l)
						assert.Equal(t, chunk.v, string(buf))
					} else {
						break
					}
				} else {
					assert.EqualError(t, err, tt.expecterr)
				}
			}
		})
	}
}
