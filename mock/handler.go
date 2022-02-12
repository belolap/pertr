// Code generated by MockGen. DO NOT EDIT.
// Source: net/http (interfaces: Handler)

// Package mock is a generated GoMock package.
package mock

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Handler is a mock of Handler interface.
type Handler struct {
	ctrl     *gomock.Controller
	recorder *HandlerMockRecorder
}

// HandlerMockRecorder is the mock recorder for Handler.
type HandlerMockRecorder struct {
	mock *Handler
}

// NewHandler creates a new mock instance.
func NewHandler(ctrl *gomock.Controller) *Handler {
	mock := &Handler{ctrl: ctrl}
	mock.recorder = &HandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Handler) EXPECT() *HandlerMockRecorder {
	return m.recorder
}

// ServeHTTP mocks base method.
func (m *Handler) ServeHTTP(arg0 http.ResponseWriter, arg1 *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServeHTTP", arg0, arg1)
}

// ServeHTTP indicates an expected call of ServeHTTP.
func (mr *HandlerMockRecorder) ServeHTTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeHTTP", reflect.TypeOf((*Handler)(nil).ServeHTTP), arg0, arg1)
}