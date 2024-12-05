// Code generated by mockery v2.46.2. DO NOT EDIT.

package log

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockHandler is an autogenerated mock type for the Handler type
type MockHandler struct {
	mock.Mock
}

// handle provides a mock function with given fields: ctx, level, log, opts
func (_m *MockHandler) handle(ctx context.Context, level LogLevel, log any, opts ...Option) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, level, log)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// NewMockHandler creates a new instance of MockHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHandler {
	mock := &MockHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
