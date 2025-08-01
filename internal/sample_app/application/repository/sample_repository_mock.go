// Code generated by MockGen. DO NOT EDIT.
// Source: sample_repository_interface.go
//
// Generated by this command:
//
//	mockgen -source=sample_repository_interface.go -destination=sample_repository_mock.go -package=repository -mock_names=ISampleRepository=MockSampleRepository
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	transaction "modern-dev-env-app-sample/internal/sample_app/application/repository/transaction"
	sample "modern-dev-env-app-sample/internal/sample_app/domain/entity/sample"
	value "modern-dev-env-app-sample/internal/sample_app/domain/value"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSampleRepository is a mock of ISampleRepository interface.
type MockSampleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSampleRepositoryMockRecorder
	isgomock struct{}
}

// MockSampleRepositoryMockRecorder is the mock recorder for MockSampleRepository.
type MockSampleRepositoryMockRecorder struct {
	mock *MockSampleRepository
}

// NewMockSampleRepository creates a new mock instance.
func NewMockSampleRepository(ctrl *gomock.Controller) *MockSampleRepository {
	mock := &MockSampleRepository{ctrl: ctrl}
	mock.recorder = &MockSampleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSampleRepository) EXPECT() *MockSampleRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSampleRepository) Delete(ctx context.Context, arg1 *sample.Sample, iTx transaction.ITransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, arg1, iTx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSampleRepositoryMockRecorder) Delete(ctx, arg1, iTx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSampleRepository)(nil).Delete), ctx, arg1, iTx)
}

// FindAll mocks base method.
func (m *MockSampleRepository) FindAll(ctx context.Context, iTx transaction.ITransaction) ([]*sample.Sample, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx, iTx)
	ret0, _ := ret[0].([]*sample.Sample)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockSampleRepositoryMockRecorder) FindAll(ctx, iTx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockSampleRepository)(nil).FindAll), ctx, iTx)
}

// FindByIDs mocks base method.
func (m *MockSampleRepository) FindByIDs(ctx context.Context, ids value.SampleIDs, iTx transaction.ITransaction) ([]*sample.Sample, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDs", ctx, ids, iTx)
	ret0, _ := ret[0].([]*sample.Sample)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByIDs indicates an expected call of FindByIDs.
func (mr *MockSampleRepositoryMockRecorder) FindByIDs(ctx, ids, iTx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDs", reflect.TypeOf((*MockSampleRepository)(nil).FindByIDs), ctx, ids, iTx)
}

// Save mocks base method.
func (m *MockSampleRepository) Save(ctx context.Context, arg1 *sample.Sample, iTx transaction.ITransaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, arg1, iTx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockSampleRepositoryMockRecorder) Save(ctx, arg1, iTx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSampleRepository)(nil).Save), ctx, arg1, iTx)
}
