// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/mocks/wes/interfaces.go

// Package wesmocks is a generated GoMock package.
package wesmocks

import (
	context "context"
	io "io"
	reflect "reflect"

	option "github.com/aws/amazon-genomics-cli/internal/pkg/wes/option"
	gomock "github.com/golang/mock/gomock"
	wes_client "github.com/rsc/wes_client"
)

// MockWesClient is a mock of WesClient interface.
type MockWesClient struct {
	ctrl     *gomock.Controller
	recorder *MockWesClientMockRecorder
}

// MockWesClientMockRecorder is the mock recorder for MockWesClient.
type MockWesClientMockRecorder struct {
	mock *MockWesClient
}

// NewMockWesClient creates a new mock instance.
func NewMockWesClient(ctrl *gomock.Controller) *MockWesClient {
	mock := &MockWesClient{ctrl: ctrl}
	mock.recorder = &MockWesClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWesClient) EXPECT() *MockWesClientMockRecorder {
	return m.recorder
}

// GetRunLog mocks base method.
func (m *MockWesClient) GetRunLog(ctx context.Context, runId string) (wes_client.RunLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunLog", ctx, runId)
	ret0, _ := ret[0].(wes_client.RunLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRunLog indicates an expected call of GetRunLog.
func (mr *MockWesClientMockRecorder) GetRunLog(ctx, runId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunLog", reflect.TypeOf((*MockWesClient)(nil).GetRunLog), ctx, runId)
}

// GetRunLogData mocks base method.
func (m *MockWesClient) GetRunLogData(ctx context.Context, runId, dataUrl string) (*io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunLogData", ctx, runId, dataUrl)
	ret0, _ := ret[0].(*io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRunLogData indicates an expected call of GetRunLogData.
func (mr *MockWesClientMockRecorder) GetRunLogData(ctx, runId, dataUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunLogData", reflect.TypeOf((*MockWesClient)(nil).GetRunLogData), ctx, runId, dataUrl)
}

// GetRunStatus mocks base method.
func (m *MockWesClient) GetRunStatus(ctx context.Context, runId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunStatus", ctx, runId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRunStatus indicates an expected call of GetRunStatus.
func (mr *MockWesClientMockRecorder) GetRunStatus(ctx, runId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunStatus", reflect.TypeOf((*MockWesClient)(nil).GetRunStatus), ctx, runId)
}

// RunWorkflow mocks base method.
func (m *MockWesClient) RunWorkflow(ctx context.Context, options ...option.Func) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RunWorkflow", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunWorkflow indicates an expected call of RunWorkflow.
func (mr *MockWesClientMockRecorder) RunWorkflow(ctx interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunWorkflow", reflect.TypeOf((*MockWesClient)(nil).RunWorkflow), varargs...)
}

// StopWorkflow mocks base method.
func (m *MockWesClient) StopWorkflow(ctx context.Context, runId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopWorkflow", ctx, runId)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopWorkflow indicates an expected call of StopWorkflow.
func (mr *MockWesClientMockRecorder) StopWorkflow(ctx, runId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopWorkflow", reflect.TypeOf((*MockWesClient)(nil).StopWorkflow), ctx, runId)
}
