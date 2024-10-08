// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Projects-WG\CLI-Project\internal\domain\interfaces\auth_service_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// IsEmailUnique mocks base method.
func (m *MockAuthService) IsEmailUnique(email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEmailUnique", email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsEmailUnique indicates an expected call of IsEmailUnique.
func (mr *MockAuthServiceMockRecorder) IsEmailUnique(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEmailUnique", reflect.TypeOf((*MockAuthService)(nil).IsEmailUnique), email)
}

// IsLeetcodeIDUnique mocks base method.
func (m *MockAuthService) IsLeetcodeIDUnique(LeetcodeID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLeetcodeIDUnique", LeetcodeID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLeetcodeIDUnique indicates an expected call of IsLeetcodeIDUnique.
func (mr *MockAuthServiceMockRecorder) IsLeetcodeIDUnique(LeetcodeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLeetcodeIDUnique", reflect.TypeOf((*MockAuthService)(nil).IsLeetcodeIDUnique), LeetcodeID)
}

// IsUsernameUnique mocks base method.
func (m *MockAuthService) IsUsernameUnique(username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUsernameUnique", username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUsernameUnique indicates an expected call of IsUsernameUnique.
func (mr *MockAuthServiceMockRecorder) IsUsernameUnique(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUsernameUnique", reflect.TypeOf((*MockAuthService)(nil).IsUsernameUnique), username)
}

// ValidateLeetcodeUsername mocks base method.
func (m *MockAuthService) ValidateLeetcodeUsername(username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateLeetcodeUsername", username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateLeetcodeUsername indicates an expected call of ValidateLeetcodeUsername.
func (mr *MockAuthServiceMockRecorder) ValidateLeetcodeUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateLeetcodeUsername", reflect.TypeOf((*MockAuthService)(nil).ValidateLeetcodeUsername), username)
}
