// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=./mock/service.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	internal "github.com/fmiskovic/cash-me-if-you-can/internal"
	account "github.com/fmiskovic/cash-me-if-you-can/internal/account"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(arg0 context.Context, arg1 *account.Account) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), arg0, arg1)
}

// Get mocks base method.
func (m *MockRepository) Get(arg0 context.Context, arg1 string) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockRepository) List(arg0 context.Context, arg1 internal.PageRequest) (internal.Page[account.Account], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(internal.Page[account.Account])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRepositoryMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepository)(nil).List), arg0, arg1)
}
