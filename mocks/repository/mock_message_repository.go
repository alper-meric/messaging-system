// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	models "github.com/alper.meric/messaging-system/models"
	mock "github.com/stretchr/testify/mock"
)

// MessageRepository is an autogenerated mock type for the MessageRepository type
type MessageRepository struct {
	mock.Mock
}

type MessageRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MessageRepository) EXPECT() *MessageRepository_Expecter {
	return &MessageRepository_Expecter{mock: &_m.Mock}
}

// AddMessage provides a mock function with given fields: message
func (_m *MessageRepository) AddMessage(message models.Message) (int, error) {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for AddMessage")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(models.Message) (int, error)); ok {
		return rf(message)
	}
	if rf, ok := ret.Get(0).(func(models.Message) int); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(models.Message) error); ok {
		r1 = rf(message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MessageRepository_AddMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddMessage'
type MessageRepository_AddMessage_Call struct {
	*mock.Call
}

// AddMessage is a helper method to define mock.On call
//   - message models.Message
func (_e *MessageRepository_Expecter) AddMessage(message interface{}) *MessageRepository_AddMessage_Call {
	return &MessageRepository_AddMessage_Call{Call: _e.mock.On("AddMessage", message)}
}

func (_c *MessageRepository_AddMessage_Call) Run(run func(message models.Message)) *MessageRepository_AddMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(models.Message))
	})
	return _c
}

func (_c *MessageRepository_AddMessage_Call) Return(_a0 int, _a1 error) *MessageRepository_AddMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MessageRepository_AddMessage_Call) RunAndReturn(run func(models.Message) (int, error)) *MessageRepository_AddMessage_Call {
	_c.Call.Return(run)
	return _c
}

// GetSentMessages provides a mock function with given fields: page, limit
func (_m *MessageRepository) GetSentMessages(page int, limit int) ([]models.Message, int, error) {
	ret := _m.Called(page, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetSentMessages")
	}

	var r0 []models.Message
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(int, int) ([]models.Message, int, error)); ok {
		return rf(page, limit)
	}
	if rf, ok := ret.Get(0).(func(int, int) []models.Message); ok {
		r0 = rf(page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) int); ok {
		r1 = rf(page, limit)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(int, int) error); ok {
		r2 = rf(page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MessageRepository_GetSentMessages_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSentMessages'
type MessageRepository_GetSentMessages_Call struct {
	*mock.Call
}

// GetSentMessages is a helper method to define mock.On call
//   - page int
//   - limit int
func (_e *MessageRepository_Expecter) GetSentMessages(page interface{}, limit interface{}) *MessageRepository_GetSentMessages_Call {
	return &MessageRepository_GetSentMessages_Call{Call: _e.mock.On("GetSentMessages", page, limit)}
}

func (_c *MessageRepository_GetSentMessages_Call) Run(run func(page int, limit int)) *MessageRepository_GetSentMessages_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(int))
	})
	return _c
}

func (_c *MessageRepository_GetSentMessages_Call) Return(_a0 []models.Message, _a1 int, _a2 error) *MessageRepository_GetSentMessages_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MessageRepository_GetSentMessages_Call) RunAndReturn(run func(int, int) ([]models.Message, int, error)) *MessageRepository_GetSentMessages_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnsentMessages provides a mock function with given fields: limit
func (_m *MessageRepository) GetUnsentMessages(limit int) ([]models.Message, error) {
	ret := _m.Called(limit)

	if len(ret) == 0 {
		panic("no return value specified for GetUnsentMessages")
	}

	var r0 []models.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]models.Message, error)); ok {
		return rf(limit)
	}
	if rf, ok := ret.Get(0).(func(int) []models.Message); ok {
		r0 = rf(limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MessageRepository_GetUnsentMessages_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnsentMessages'
type MessageRepository_GetUnsentMessages_Call struct {
	*mock.Call
}

// GetUnsentMessages is a helper method to define mock.On call
//   - limit int
func (_e *MessageRepository_Expecter) GetUnsentMessages(limit interface{}) *MessageRepository_GetUnsentMessages_Call {
	return &MessageRepository_GetUnsentMessages_Call{Call: _e.mock.On("GetUnsentMessages", limit)}
}

func (_c *MessageRepository_GetUnsentMessages_Call) Run(run func(limit int)) *MessageRepository_GetUnsentMessages_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MessageRepository_GetUnsentMessages_Call) Return(_a0 []models.Message, _a1 error) *MessageRepository_GetUnsentMessages_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MessageRepository_GetUnsentMessages_Call) RunAndReturn(run func(int) ([]models.Message, error)) *MessageRepository_GetUnsentMessages_Call {
	_c.Call.Return(run)
	return _c
}

// MarkMessageAsSent provides a mock function with given fields: id, externalMsgID
func (_m *MessageRepository) MarkMessageAsSent(id int, externalMsgID string) error {
	ret := _m.Called(id, externalMsgID)

	if len(ret) == 0 {
		panic("no return value specified for MarkMessageAsSent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(id, externalMsgID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MessageRepository_MarkMessageAsSent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarkMessageAsSent'
type MessageRepository_MarkMessageAsSent_Call struct {
	*mock.Call
}

// MarkMessageAsSent is a helper method to define mock.On call
//   - id int
//   - externalMsgID string
func (_e *MessageRepository_Expecter) MarkMessageAsSent(id interface{}, externalMsgID interface{}) *MessageRepository_MarkMessageAsSent_Call {
	return &MessageRepository_MarkMessageAsSent_Call{Call: _e.mock.On("MarkMessageAsSent", id, externalMsgID)}
}

func (_c *MessageRepository_MarkMessageAsSent_Call) Run(run func(id int, externalMsgID string)) *MessageRepository_MarkMessageAsSent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(string))
	})
	return _c
}

func (_c *MessageRepository_MarkMessageAsSent_Call) Return(_a0 error) *MessageRepository_MarkMessageAsSent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MessageRepository_MarkMessageAsSent_Call) RunAndReturn(run func(int, string) error) *MessageRepository_MarkMessageAsSent_Call {
	_c.Call.Return(run)
	return _c
}

// NewMessageRepository creates a new instance of MessageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMessageRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MessageRepository {
	mock := &MessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
