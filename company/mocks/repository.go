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

	model "github.com/bopoh24/ma_1/company/internal/model"
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

// CompanyActivateDeactivate mocks base method.
func (m *MockRepository) CompanyActivateDeactivate(ctx context.Context, id int64, active bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyActivateDeactivate", ctx, id, active)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompanyActivateDeactivate indicates an expected call of CompanyActivateDeactivate.
func (mr *MockRepositoryMockRecorder) CompanyActivateDeactivate(ctx, id, active any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyActivateDeactivate", reflect.TypeOf((*MockRepository)(nil).CompanyActivateDeactivate), ctx, id, active)
}

// CompanyByID mocks base method.
func (m *MockRepository) CompanyByID(ctx context.Context, id int64) (model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyByID", ctx, id)
	ret0, _ := ret[0].(model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompanyByID indicates an expected call of CompanyByID.
func (mr *MockRepositoryMockRecorder) CompanyByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyByID", reflect.TypeOf((*MockRepository)(nil).CompanyByID), ctx, id)
}

// CompanyCreate mocks base method.
func (m *MockRepository) CompanyCreate(ctx context.Context, userId, email, firstName, lastName string, company model.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyCreate", ctx, userId, email, firstName, lastName, company)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompanyCreate indicates an expected call of CompanyCreate.
func (mr *MockRepositoryMockRecorder) CompanyCreate(ctx, userId, email, firstName, lastName, company any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyCreate", reflect.TypeOf((*MockRepository)(nil).CompanyCreate), ctx, userId, email, firstName, lastName, company)
}

// CompanyManagers mocks base method.
func (m *MockRepository) CompanyManagers(ctx context.Context, companyID int64) ([]model.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyManagers", ctx, companyID)
	ret0, _ := ret[0].([]model.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompanyManagers indicates an expected call of CompanyManagers.
func (mr *MockRepositoryMockRecorder) CompanyManagers(ctx, companyID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyManagers", reflect.TypeOf((*MockRepository)(nil).CompanyManagers), ctx, companyID)
}

// CompanyUpdate mocks base method.
func (m *MockRepository) CompanyUpdate(ctx context.Context, company model.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyUpdate", ctx, company)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompanyUpdate indicates an expected call of CompanyUpdate.
func (mr *MockRepositoryMockRecorder) CompanyUpdate(ctx, company any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyUpdate", reflect.TypeOf((*MockRepository)(nil).CompanyUpdate), ctx, company)
}

// CompanyUpdateLocation mocks base method.
func (m *MockRepository) CompanyUpdateLocation(ctx context.Context, id int64, lat, lng float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyUpdateLocation", ctx, id, lat, lng)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompanyUpdateLocation indicates an expected call of CompanyUpdateLocation.
func (mr *MockRepositoryMockRecorder) CompanyUpdateLocation(ctx, id, lat, lng any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyUpdateLocation", reflect.TypeOf((*MockRepository)(nil).CompanyUpdateLocation), ctx, id, lat, lng)
}

// CompanyUpdateLogo mocks base method.
func (m *MockRepository) CompanyUpdateLogo(ctx context.Context, id int64, logo string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompanyUpdateLogo", ctx, id, logo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompanyUpdateLogo indicates an expected call of CompanyUpdateLogo.
func (mr *MockRepositoryMockRecorder) CompanyUpdateLogo(ctx, id, logo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompanyUpdateLogo", reflect.TypeOf((*MockRepository)(nil).CompanyUpdateLogo), ctx, id, logo)
}

// ManagerActivateDeactivate mocks base method.
func (m *MockRepository) ManagerActivateDeactivate(ctx context.Context, id int64, active bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerActivateDeactivate", ctx, id, active)
	ret0, _ := ret[0].(error)
	return ret0
}

// ManagerActivateDeactivate indicates an expected call of ManagerActivateDeactivate.
func (mr *MockRepositoryMockRecorder) ManagerActivateDeactivate(ctx, id, active any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerActivateDeactivate", reflect.TypeOf((*MockRepository)(nil).ManagerActivateDeactivate), ctx, id, active)
}

// ManagerByEmail mocks base method.
func (m *MockRepository) ManagerByEmail(ctx context.Context, email string) (model.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerByEmail", ctx, email)
	ret0, _ := ret[0].(model.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ManagerByEmail indicates an expected call of ManagerByEmail.
func (mr *MockRepositoryMockRecorder) ManagerByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerByEmail", reflect.TypeOf((*MockRepository)(nil).ManagerByEmail), ctx, email)
}

// ManagerByID mocks base method.
func (m *MockRepository) ManagerByID(ctx context.Context, id int64) (model.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerByID", ctx, id)
	ret0, _ := ret[0].(model.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ManagerByID indicates an expected call of ManagerByID.
func (mr *MockRepositoryMockRecorder) ManagerByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerByID", reflect.TypeOf((*MockRepository)(nil).ManagerByID), ctx, id)
}

// ManagerByUserID mocks base method.
func (m *MockRepository) ManagerByUserID(ctx context.Context, userId string) (model.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerByUserID", ctx, userId)
	ret0, _ := ret[0].(model.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ManagerByUserID indicates an expected call of ManagerByUserID.
func (mr *MockRepositoryMockRecorder) ManagerByUserID(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerByUserID", reflect.TypeOf((*MockRepository)(nil).ManagerByUserID), ctx, userId)
}

// ManagerDelete mocks base method.
func (m *MockRepository) ManagerDelete(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerDelete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ManagerDelete indicates an expected call of ManagerDelete.
func (mr *MockRepositoryMockRecorder) ManagerDelete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerDelete", reflect.TypeOf((*MockRepository)(nil).ManagerDelete), ctx, id)
}

// ManagerInvite mocks base method.
func (m *MockRepository) ManagerInvite(ctx context.Context, companyID int64, email string, role model.MangerRole) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerInvite", ctx, companyID, email, role)
	ret0, _ := ret[0].(error)
	return ret0
}

// ManagerInvite indicates an expected call of ManagerInvite.
func (mr *MockRepositoryMockRecorder) ManagerInvite(ctx, companyID, email, role any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerInvite", reflect.TypeOf((*MockRepository)(nil).ManagerInvite), ctx, companyID, email, role)
}

// ManagerRole mocks base method.
func (m *MockRepository) ManagerRole(ctx context.Context, companyId int64, userId string) (model.MangerRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerRole", ctx, companyId, userId)
	ret0, _ := ret[0].(model.MangerRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ManagerRole indicates an expected call of ManagerRole.
func (mr *MockRepositoryMockRecorder) ManagerRole(ctx, companyId, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerRole", reflect.TypeOf((*MockRepository)(nil).ManagerRole), ctx, companyId, userId)
}

// ManagerSetRole mocks base method.
func (m *MockRepository) ManagerSetRole(ctx context.Context, id int64, role model.MangerRole) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ManagerSetRole", ctx, id, role)
	ret0, _ := ret[0].(error)
	return ret0
}

// ManagerSetRole indicates an expected call of ManagerSetRole.
func (mr *MockRepositoryMockRecorder) ManagerSetRole(ctx, id, role any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ManagerSetRole", reflect.TypeOf((*MockRepository)(nil).ManagerSetRole), ctx, id, role)
}
