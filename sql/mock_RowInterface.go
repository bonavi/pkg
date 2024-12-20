// Code generated by mockery v2.46.2. DO NOT EDIT.

package sql

import mock "github.com/stretchr/testify/mock"

// MockRowInterface is an autogenerated mock type for the RowInterface type
type MockRowInterface struct {
	mock.Mock
}

// MapScan provides a mock function with given fields: dest
func (_m *MockRowInterface) MapScan(dest map[string]any) error {
	ret := _m.Called(dest)

	if len(ret) == 0 {
		panic("no return value specified for MapScan")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]any) error); ok {
		r0 = rf(dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Scan provides a mock function with given fields: dest
func (_m *MockRowInterface) Scan(dest ...any) error {
	var _ca []interface{}
	_ca = append(_ca, dest...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Scan")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...any) error); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SliceScan provides a mock function with given fields:
func (_m *MockRowInterface) SliceScan() ([]any, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SliceScan")
	}

	var r0 []any
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]any, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []any); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]any)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StructScan provides a mock function with given fields: dest
func (_m *MockRowInterface) StructScan(dest any) error {
	ret := _m.Called(dest)

	if len(ret) == 0 {
		panic("no return value specified for StructScan")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(any) error); ok {
		r0 = rf(dest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockRowInterface creates a new instance of MockRowInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRowInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRowInterface {
	mock := &MockRowInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
