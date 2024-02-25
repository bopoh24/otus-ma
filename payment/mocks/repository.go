// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source service.go -destination ../../mocks/repository.go -package mock Repository
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

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

// Balance mocks base method.
func (m *MockRepository) Balance(ctx context.Context, customerID string) (float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Balance", ctx, customerID)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Balance indicates an expected call of Balance.
func (mr *MockRepositoryMockRecorder) Balance(ctx, customerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Balance", reflect.TypeOf((*MockRepository)(nil).Balance), ctx, customerID)
}

// Close mocks base method.
func (m *MockRepository) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close), ctx)
}

// CreateAccount mocks base method.
func (m *MockRepository) CreateAccount(ctx context.Context, customerID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, customerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockRepositoryMockRecorder) CreateAccount(ctx, customerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockRepository)(nil).CreateAccount), ctx, customerID)
}

// PaymentCancel mocks base method.
func (m *MockRepository) PaymentCancel(ctx context.Context, orderId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaymentCancel", ctx, orderId)
	ret0, _ := ret[0].(error)
	return ret0
}

// PaymentCancel indicates an expected call of PaymentCancel.
func (mr *MockRepositoryMockRecorder) PaymentCancel(ctx, orderId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaymentCancel", reflect.TypeOf((*MockRepository)(nil).PaymentCancel), ctx, orderId)
}

// PaymentMake mocks base method.
func (m *MockRepository) PaymentMake(ctx context.Context, orderId int64, customerID string, amount float32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaymentMake", ctx, orderId, customerID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// PaymentMake indicates an expected call of PaymentMake.
func (mr *MockRepositoryMockRecorder) PaymentMake(ctx, orderId, customerID, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaymentMake", reflect.TypeOf((*MockRepository)(nil).PaymentMake), ctx, orderId, customerID, amount)
}

// TopUp mocks base method.
func (m *MockRepository) TopUp(ctx context.Context, customerID string, amount float32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TopUp", ctx, customerID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// TopUp indicates an expected call of TopUp.
func (mr *MockRepositoryMockRecorder) TopUp(ctx, customerID, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TopUp", reflect.TypeOf((*MockRepository)(nil).TopUp), ctx, customerID, amount)
}
