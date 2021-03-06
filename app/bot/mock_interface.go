// Code generated by mockery v2.10.4. DO NOT EDIT.

package bot

import mock "github.com/stretchr/testify/mock"

// MockInterface is an autogenerated mock type for the Interface type
type MockInterface struct {
	mock.Mock
}

// Help provides a mock function with given fields:
func (_m *MockInterface) Help() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OnMessage provides a mock function with given fields: msg
func (_m *MockInterface) OnMessage(msg Message) Response {
	ret := _m.Called(msg)

	var r0 Response
	if rf, ok := ret.Get(0).(func(Message) Response); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Get(0).(Response)
	}

	return r0
}

// ReactOn provides a mock function with given fields:
func (_m *MockInterface) ReactOn() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
