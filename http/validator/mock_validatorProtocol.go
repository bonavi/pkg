// Code generated by mockery v2.46.2. DO NOT EDIT.

package middleware

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockValidatorProtocol is an autogenerated mock type for the validatorProtocol type
type mockValidatorProtocol struct {
	mock.Mock
}

// Validate provides a mock function with given fields: ctx
func (_m *mockValidatorProtocol) Validate(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockValidatorProtocol creates a new instance of mockValidatorProtocol. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockValidatorProtocol(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockValidatorProtocol {
	mock := &mockValidatorProtocol{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
