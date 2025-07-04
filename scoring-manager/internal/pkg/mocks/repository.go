// Code generated by MockGen. DO NOT EDIT.
// Source: scoring-manager/internal/client/repository/repository.go
//
// Generated by this command:
//
//	mockgen -source=scoring-manager/internal/client/repository/repository.go -destination=scoring-manager/internal/pkg/mocks/repository.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	repository "github.com/KlassnayaAfrodita/github-user-score/scoring-manager/internal/client/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockScoringRepositoryInterface is a mock of ScoringRepositoryInterface interface.
type MockScoringRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockScoringRepositoryInterfaceMockRecorder
	isgomock struct{}
}

// MockScoringRepositoryInterfaceMockRecorder is the mock recorder for MockScoringRepositoryInterface.
type MockScoringRepositoryInterfaceMockRecorder struct {
	mock *MockScoringRepositoryInterface
}

// NewMockScoringRepositoryInterface creates a new mock instance.
func NewMockScoringRepositoryInterface(ctrl *gomock.Controller) *MockScoringRepositoryInterface {
	mock := &MockScoringRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockScoringRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScoringRepositoryInterface) EXPECT() *MockScoringRepositoryInterfaceMockRecorder {
	return m.recorder
}

// CreateScoringApplication mocks base method.
func (m *MockScoringRepositoryInterface) CreateScoringApplication(ctx context.Context, app repository.ScoringApplication) (repository.ScoringApplication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateScoringApplication", ctx, app)
	ret0, _ := ret[0].(repository.ScoringApplication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateScoringApplication indicates an expected call of CreateScoringApplication.
func (mr *MockScoringRepositoryInterfaceMockRecorder) CreateScoringApplication(ctx, app any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateScoringApplication", reflect.TypeOf((*MockScoringRepositoryInterface)(nil).CreateScoringApplication), ctx, app)
}

// GetExpiredApplications mocks base method.
func (m *MockScoringRepositoryInterface) GetExpiredApplications(ctx context.Context, maxAgeMinutes int) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpiredApplications", ctx, maxAgeMinutes)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExpiredApplications indicates an expected call of GetExpiredApplications.
func (mr *MockScoringRepositoryInterfaceMockRecorder) GetExpiredApplications(ctx, maxAgeMinutes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpiredApplications", reflect.TypeOf((*MockScoringRepositoryInterface)(nil).GetExpiredApplications), ctx, maxAgeMinutes)
}

// GetScoringApplicationByID mocks base method.
func (m *MockScoringRepositoryInterface) GetScoringApplicationByID(ctx context.Context, appID int64) (repository.ScoringApplication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScoringApplicationByID", ctx, appID)
	ret0, _ := ret[0].(repository.ScoringApplication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScoringApplicationByID indicates an expected call of GetScoringApplicationByID.
func (mr *MockScoringRepositoryInterfaceMockRecorder) GetScoringApplicationByID(ctx, appID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScoringApplicationByID", reflect.TypeOf((*MockScoringRepositoryInterface)(nil).GetScoringApplicationByID), ctx, appID)
}

// MarkExpiredApplications mocks base method.
func (m *MockScoringRepositoryInterface) MarkExpiredApplications(ctx context.Context, appIDs []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkExpiredApplications", ctx, appIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkExpiredApplications indicates an expected call of MarkExpiredApplications.
func (mr *MockScoringRepositoryInterfaceMockRecorder) MarkExpiredApplications(ctx, appIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkExpiredApplications", reflect.TypeOf((*MockScoringRepositoryInterface)(nil).MarkExpiredApplications), ctx, appIDs)
}

// SaveScoringApplicationResult mocks base method.
func (m *MockScoringRepositoryInterface) SaveScoringApplicationResult(ctx context.Context, app repository.ScoringApplication) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveScoringApplicationResult", ctx, app)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveScoringApplicationResult indicates an expected call of SaveScoringApplicationResult.
func (mr *MockScoringRepositoryInterfaceMockRecorder) SaveScoringApplicationResult(ctx, app any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveScoringApplicationResult", reflect.TypeOf((*MockScoringRepositoryInterface)(nil).SaveScoringApplicationResult), ctx, app)
}

// UpdateScoringApplicationStatus mocks base method.
func (m *MockScoringRepositoryInterface) UpdateScoringApplicationStatus(ctx context.Context, appID int64, status repository.ScoringStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateScoringApplicationStatus", ctx, appID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateScoringApplicationStatus indicates an expected call of UpdateScoringApplicationStatus.
func (mr *MockScoringRepositoryInterfaceMockRecorder) UpdateScoringApplicationStatus(ctx, appID, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateScoringApplicationStatus", reflect.TypeOf((*MockScoringRepositoryInterface)(nil).UpdateScoringApplicationStatus), ctx, appID, status)
}
