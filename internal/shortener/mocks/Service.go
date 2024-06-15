// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CrateAlias provides a mock function with given fields: url
func (_m *Service) CrateAlias(url string) (string, error) {
	ret := _m.Called(url)

	if len(ret) == 0 {
		panic("no return value specified for SaveAlias")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(url)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrigURL provides a mock function with given fields: alias
func (_m *Service) GetOrigURL(alias string) (string, error) {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for GetOrigURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
