// Code generated by MockGen. DO NOT EDIT.
// Source: net/http (interfaces: RoundTripper)

// Package mock is a generated GoMock package.
package mock

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Transport is a mock of RoundTripper interface.
type Transport struct {
	ctrl     *gomock.Controller
	recorder *TransportMockRecorder
}

// TransportMockRecorder is the mock recorder for Transport.
type TransportMockRecorder struct {
	mock *Transport
}

// NewTransport creates a new mock instance.
func NewTransport(ctrl *gomock.Controller) *Transport {
	mock := &Transport{ctrl: ctrl}
	mock.recorder = &TransportMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Transport) EXPECT() *TransportMockRecorder {
	return m.recorder
}

// RoundTrip mocks base method.
func (m *Transport) RoundTrip(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoundTrip", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RoundTrip indicates an expected call of RoundTrip.
func (mr *TransportMockRecorder) RoundTrip(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoundTrip", reflect.TypeOf((*Transport)(nil).RoundTrip), arg0)
}
