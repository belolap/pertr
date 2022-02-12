package pertr

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/belolap/pertr/match"
	"github.com/belolap/pertr/mock"
)

// testExpectChunk is expected return pair from body.Read().
type testExpectChunk struct {
	num   int
	value string
}

func TestBody(t *testing.T) {
	tests := []struct {
		name string

		// Request settings.
		method  string        // Request method.
		timeout time.Duration // Timeout for download retries context.
		sizes   []int         // Body read buf sizes.

		// Set of handler functions for seq server response.
		handler func(*mock.Handler)

		// Setup body mock function.
		body func(*mock.Body)

		// Expectations.
		expect    []testExpectChunk // expect .Read() returns
		expectErr string
	}{
		{
			name:   "simple request",
			method: "GET",
			sizes:  []int{3},
			handler: func(h *mock.Handler) {
				h.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
			},
			body: func(b *mock.Body) {
				b.EXPECT().Read(match.SliceSize(3)).Return(3, nil).Do(func(p []byte) {
					copy(p, []byte("abc"))
				})
				b.EXPECT().Close()
			},
			expect: []testExpectChunk{{num: 3, value: "abc"}},
		},
		{
			name:   "immediately error",
			method: "GET",
			sizes:  []int{1},
			handler: func(h *mock.Handler) {
				gomock.InOrder(
					h.EXPECT().ServeHTTP(gomock.Any(), gomock.Any()).Do(func(w http.ResponseWriter, r *http.Request) {
						_, _ = io.WriteString(w, "abc")
					}),
					h.EXPECT().ServeHTTP(gomock.Any(), gomock.Any()).Do(func(w http.ResponseWriter, r *http.Request) {
						_, _ = io.WriteString(w, "bc")
					}),
				)

			},
			body: func(b *mock.Body) {
				b.EXPECT().Read(match.SliceSize(1)).Return(0, errors.New("some err"))
				b.EXPECT().Close()
			},
			expectErr: "some err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize gomock.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Start test http server.
			handler := mock.NewHandler(ctrl)
			if tt.handler != nil {
				tt.handler(handler)
			}

			ts := httptest.NewServer(handler)
			defer ts.Close()

			// Create test transport with body read error emulation.
			body := mock.NewBody(ctrl)
			if tt.body != nil {
				tt.body(body)
			}
			trans := mock.NewTransport(ctrl)
			trans.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				rsp, err := http.DefaultTransport.RoundTrip(req)
				rsp.Body = body
				return rsp, err
			})

			// Set request context with download timeout.
			ctx := context.Background()
			if tt.timeout != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.timeout)
				defer cancel()
			}

			// Create our wrapper.
			w := New(WithTransport(trans), WithContext(ctx))

			// Create client.
			cli := &http.Client{Transport: w}

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
			for i, size := range tt.sizes {
				buf := make([]byte, size)
				var n int
				n, err = rsp.Body.Read(buf)
				if err == nil || errors.Is(err, io.EOF) {
					if i < len(tt.expect) {
						if n != tt.expect[i].num {
							t.Errorf("Invalid number of bytes returned, want %d, got %d.", tt.expect[i].num, n)
						}
						if string(buf) != tt.expect[i].value {
							t.Errorf("Invalid downloaded body, want `%s\", got `%s\".", tt.expect[i].value, string(buf))
						}
					}
					err = nil
				} else {
					break
				}
			}
			if tt.expectErr == "" {
				if err != nil {
					t.Errorf("Unexpected error: %s.", err.Error())
				}
			} else {
				if err == nil {
					t.Error("Expect error, but got nil.")
				} else if err.Error() != tt.expectErr {
					t.Errorf("Invalid error recieved: want `%s\", got `%s\".", tt.expectErr, err.Error())
				}
			}
		})
	}
}
