// Code generated by mockery v2.46.2. DO NOT EDIT.

package chain

import mock "github.com/stretchr/testify/mock"

// MockOption is an autogenerated mock type for the Option type
type MockOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *MockOption) Execute(_a0 *Chain) {
	_m.Called(_a0)
}

// NewMockOption creates a new instance of MockOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOption {
	mock := &MockOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
