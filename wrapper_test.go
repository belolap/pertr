package pertr

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Micro mock for transport.
type testFakeTransport struct{}

// RoundTrip
func (t *testFakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case "GET":
		return &http.Response{Request: req}, nil
	default:
		return nil, errors.New("some err")
	}
}

func TestWrapper(t *testing.T) {
	tests := []struct {
		up        *http.Transport
		name      string
		method    string
		expectErr string
	}{
		{
			name:   "positive req",
			method: "GET",
		},
		{
			name:      "negative req",
			method:    "POST",
			expectErr: "Post \"blah\": some err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &http.Client{}
			cli.Transport = New(WithClient(cli), WithTransport(&testFakeTransport{}), WithContext(context.TODO()))

			req, err := http.NewRequest(tt.method, "blah", http.NoBody)
			if err != nil {
				panic(err)
			}

			rsp, err := cli.Do(req)
			if tt.expectErr == "" {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.method, rsp.Request.Method)
				}
			} else {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectErr)
				}
			}
		})
	}
}
