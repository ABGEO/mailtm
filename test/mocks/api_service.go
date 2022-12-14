// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	dto "github.com/abgeo/mailtm/pkg/dto"
	mock "github.com/stretchr/testify/mock"
)

// APIServiceInterface is an autogenerated mocks type for the APIServiceInterface type
type APIServiceInterface struct {
	mock.Mock
}

// CreateAccount provides a mocks function with given fields: data
func (_m *APIServiceInterface) CreateAccount(data dto.AccountWrite) (*dto.Account, error) {
	ret := _m.Called(data)

	var r0 *dto.Account
	if rf, ok := ret.Get(0).(func(dto.AccountWrite) *dto.Account); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dto.AccountWrite) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DownloadMessageAttachment provides a mocks function with given fields: messageID, attachmentID, path
func (_m *APIServiceInterface) DownloadMessageAttachment(messageID string, attachmentID string, path string) error {
	ret := _m.Called(messageID, attachmentID, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(messageID, attachmentID, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccount provides a mocks function with given fields: id
func (_m *APIServiceInterface) GetAccount(id string) (*dto.Account, error) {
	ret := _m.Called(id)

	var r0 *dto.Account
	if rf, ok := ret.Get(0).(func(string) *dto.Account); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCurrentAccount provides a mocks function with given fields:
func (_m *APIServiceInterface) GetCurrentAccount() (*dto.Account, error) {
	ret := _m.Called()

	var r0 *dto.Account
	if rf, ok := ret.Get(0).(func() *dto.Account); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDomain provides a mocks function with given fields: id
func (_m *APIServiceInterface) GetDomain(id string) (*dto.Domain, error) {
	ret := _m.Called(id)

	var r0 *dto.Domain
	if rf, ok := ret.Get(0).(func(string) *dto.Domain); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.Domain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDomains provides a mocks function with given fields:
func (_m *APIServiceInterface) GetDomains() ([]dto.Domain, error) {
	ret := _m.Called()

	var r0 []dto.Domain
	if rf, ok := ret.Get(0).(func() []dto.Domain); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.Domain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessage provides a mocks function with given fields: id
func (_m *APIServiceInterface) GetMessage(id string) (*dto.Message, error) {
	ret := _m.Called(id)

	var r0 *dto.Message
	if rf, ok := ret.Get(0).(func(string) *dto.Message); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessages provides a mocks function with given fields:
func (_m *APIServiceInterface) GetMessages() (dto.Messages, error) {
	ret := _m.Called()

	var r0 dto.Messages
	if rf, ok := ret.Get(0).(func() dto.Messages); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(dto.Messages)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSource provides a mocks function with given fields: id
func (_m *APIServiceInterface) GetSource(id string) (*dto.Source, error) {
	ret := _m.Called(id)

	var r0 *dto.Source
	if rf, ok := ret.Get(0).(func(string) *dto.Source); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.Source)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetToken provides a mocks function with given fields: credentials
func (_m *APIServiceInterface) GetToken(credentials dto.Credentials) (*dto.Token, error) {
	ret := _m.Called(credentials)

	var r0 *dto.Token
	if rf, ok := ret.Get(0).(func(dto.Credentials) *dto.Token); ok {
		r0 = rf(credentials)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dto.Credentials) error); ok {
		r1 = rf(credentials)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveAccount provides a mocks function with given fields: id
func (_m *APIServiceInterface) RemoveAccount(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveMessage provides a mocks function with given fields: id
func (_m *APIServiceInterface) RemoveMessage(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetToken provides a mocks function with given fields: token
func (_m *APIServiceInterface) SetToken(token *dto.Token) {
	_m.Called(token)
}

// UpdateMessage provides a mocks function with given fields: id, data
func (_m *APIServiceInterface) UpdateMessage(id string, data dto.MessageWrite) error {
	ret := _m.Called(id, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, dto.MessageWrite) error); ok {
		r0 = rf(id, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAPIServiceInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewAPIServiceInterface creates a new instance of APIServiceInterface. It also registers a testing interface on the mocks and a cleanup function to assert the mocks expectations.
func NewAPIServiceInterface(t mockConstructorTestingTNewAPIServiceInterface) *APIServiceInterface {
	mock := &APIServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
