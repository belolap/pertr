package pertr

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/belolap/pertr/mock"
)

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
			// Initialize gomock.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create custom transport.
			trans := mock.NewTransport(ctrl)
			trans.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
				switch req.Method {
				case "GET":
					return &http.Response{Request: req}, nil
				default:
					return nil, errors.New("some err")
				}
			})

			// Creating client.
			cli := &http.Client{
				Transport: New(WithTransport(trans), WithContext(context.TODO())),
			}

			// Create request.
			req, err := http.NewRequest(tt.method, "blah", http.NoBody)
			if err != nil {
				panic(err)
			}

			// Execute.
			rsp, err := cli.Do(req)
			if tt.expectErr == "" {
				if err != nil {
					t.Errorf("Unexpected error: %s.", err.Error())
				} else if rsp.Request.Method != tt.method {
					t.Errorf("Invalid response recieved.")
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
