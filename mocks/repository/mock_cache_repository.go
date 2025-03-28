// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// CacheRepository is an autogenerated mock type for the CacheRepository type
type CacheRepository struct {
	mock.Mock
}

type CacheRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *CacheRepository) EXPECT() *CacheRepository_Expecter {
	return &CacheRepository_Expecter{mock: &_m.Mock}
}

// CacheMessageID provides a mock function with given fields: messageID, sentAt
func (_m *CacheRepository) CacheMessageID(messageID string, sentAt time.Time) error {
	ret := _m.Called(messageID, sentAt)

	if len(ret) == 0 {
		panic("no return value specified for CacheMessageID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Time) error); ok {
		r0 = rf(messageID, sentAt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CacheRepository_CacheMessageID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CacheMessageID'
type CacheRepository_CacheMessageID_Call struct {
	*mock.Call
}

// CacheMessageID is a helper method to define mock.On call
//   - messageID string
//   - sentAt time.Time
func (_e *CacheRepository_Expecter) CacheMessageID(messageID interface{}, sentAt interface{}) *CacheRepository_CacheMessageID_Call {
	return &CacheRepository_CacheMessageID_Call{Call: _e.mock.On("CacheMessageID", messageID, sentAt)}
}

func (_c *CacheRepository_CacheMessageID_Call) Run(run func(messageID string, sentAt time.Time)) *CacheRepository_CacheMessageID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(time.Time))
	})
	return _c
}

func (_c *CacheRepository_CacheMessageID_Call) Return(_a0 error) *CacheRepository_CacheMessageID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheRepository_CacheMessageID_Call) RunAndReturn(run func(string, time.Time) error) *CacheRepository_CacheMessageID_Call {
	_c.Call.Return(run)
	return _c
}

// GetCachedMessage provides a mock function with given fields: messageID
func (_m *CacheRepository) GetCachedMessage(messageID string) (time.Time, error) {
	ret := _m.Called(messageID)

	if len(ret) == 0 {
		panic("no return value specified for GetCachedMessage")
	}

	var r0 time.Time
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (time.Time, error)); ok {
		return rf(messageID)
	}
	if rf, ok := ret.Get(0).(func(string) time.Time); ok {
		r0 = rf(messageID)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(messageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CacheRepository_GetCachedMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCachedMessage'
type CacheRepository_GetCachedMessage_Call struct {
	*mock.Call
}

// GetCachedMessage is a helper method to define mock.On call
//   - messageID string
func (_e *CacheRepository_Expecter) GetCachedMessage(messageID interface{}) *CacheRepository_GetCachedMessage_Call {
	return &CacheRepository_GetCachedMessage_Call{Call: _e.mock.On("GetCachedMessage", messageID)}
}

func (_c *CacheRepository_GetCachedMessage_Call) Run(run func(messageID string)) *CacheRepository_GetCachedMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *CacheRepository_GetCachedMessage_Call) Return(_a0 time.Time, _a1 error) *CacheRepository_GetCachedMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CacheRepository_GetCachedMessage_Call) RunAndReturn(run func(string) (time.Time, error)) *CacheRepository_GetCachedMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewCacheRepository creates a new instance of CacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCacheRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CacheRepository {
	mock := &CacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
