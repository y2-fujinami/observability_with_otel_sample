// Code generated by MockGen. DO NOT EDIT.
// Source: list_samples_usecase_interface.go
//
// Generated by this command:
//
//	mockgen -source=list_samples_usecase_interface.go -destination=list_samples_usecase_mock.go -package=sample -mock_names=IListSamplesUseCase=MockListSamplesUseCase
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

// MockListSamplesUseCase is a mock of IListSamplesUseCase interface.
type MockListSamplesUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockListSamplesUseCaseMockRecorder
	isgomock struct{}
}

// MockListSamplesUseCaseMockRecorder is the mock recorder for MockListSamplesUseCase.
type MockListSamplesUseCaseMockRecorder struct {
	mock *MockListSamplesUseCase
}

// NewMockListSamplesUseCase creates a new mock instance.
func NewMockListSamplesUseCase(ctrl *gomock.Controller) *MockListSamplesUseCase {
	mock := &MockListSamplesUseCase{ctrl: ctrl}
	mock.recorder = &MockListSamplesUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockListSamplesUseCase) EXPECT() *MockListSamplesUseCaseMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockListSamplesUseCase) Run(ctx context.Context, req *sample.ListSamplesRequest) (*sample0.ListSamplesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx, req)
	ret0, _ := ret[0].(*sample0.ListSamplesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockListSamplesUseCaseMockRecorder) Run(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockListSamplesUseCase)(nil).Run), ctx, req)
}
