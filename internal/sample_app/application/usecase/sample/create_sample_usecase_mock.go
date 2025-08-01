// Code generated by MockGen. DO NOT EDIT.
// Source: create_sample_usecase_interface.go
//
// Generated by this command:
//
//	mockgen -source=create_sample_usecase_interface.go -destination=create_sample_usecase_mock.go -package=sample -mock_names=ICreateSampleUseCase=MockCreateSampleUseCase
//

// Package sample is a generated GoMock package.
package sample

import (
	context "context"
	sample "modern-dev-env-app-sample/internal/sample_app/application/request/sample"
	sample0 "modern-dev-env-app-sample/internal/sample_app/application/response/sample"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCreateSampleUseCase is a mock of ICreateSampleUseCase interface.
type MockCreateSampleUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCreateSampleUseCaseMockRecorder
	isgomock struct{}
}

// MockCreateSampleUseCaseMockRecorder is the mock recorder for MockCreateSampleUseCase.
type MockCreateSampleUseCaseMockRecorder struct {
	mock *MockCreateSampleUseCase
}

// NewMockCreateSampleUseCase creates a new mock instance.
func NewMockCreateSampleUseCase(ctrl *gomock.Controller) *MockCreateSampleUseCase {
	mock := &MockCreateSampleUseCase{ctrl: ctrl}
	mock.recorder = &MockCreateSampleUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateSampleUseCase) EXPECT() *MockCreateSampleUseCaseMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockCreateSampleUseCase) Run(ctx context.Context, req *sample.CreateSampleRequest) (*sample0.CreateSampleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx, req)
	ret0, _ := ret[0].(*sample0.CreateSampleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockCreateSampleUseCaseMockRecorder) Run(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCreateSampleUseCase)(nil).Run), ctx, req)
}
