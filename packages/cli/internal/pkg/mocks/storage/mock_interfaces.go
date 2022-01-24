// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/mocks/storage/interfaces.go

// Package storagemocks is a generated GoMock package.
package storagemocks

import (
	reflect "reflect"

	config "github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	spec "github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	gomock "github.com/golang/mock/gomock"
)

// MockProjectClient is a mock of ProjectClient interface.
type MockProjectClient struct {
	ctrl     *gomock.Controller
	recorder *MockProjectClientMockRecorder
}

// MockProjectClientMockRecorder is the mock recorder for MockProjectClient.
type MockProjectClientMockRecorder struct {
	mock *MockProjectClient
}

// NewMockProjectClient creates a new mock instance.
func NewMockProjectClient(ctrl *gomock.Controller) *MockProjectClient {
	mock := &MockProjectClient{ctrl: ctrl}
	mock.recorder = &MockProjectClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectClient) EXPECT() *MockProjectClientMockRecorder {
	return m.recorder
}

// GetLocation mocks base method.
func (m *MockProjectClient) GetLocation() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocation")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLocation indicates an expected call of GetLocation.
func (mr *MockProjectClientMockRecorder) GetLocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocation", reflect.TypeOf((*MockProjectClient)(nil).GetLocation))
}

// GetProjectName mocks base method.
func (m *MockProjectClient) GetProjectName() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectName")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectName indicates an expected call of GetProjectName.
func (mr *MockProjectClientMockRecorder) GetProjectName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectName", reflect.TypeOf((*MockProjectClient)(nil).GetProjectName))
}

// IsInitialized mocks base method.
func (m *MockProjectClient) IsInitialized() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInitialized")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsInitialized indicates an expected call of IsInitialized.
func (mr *MockProjectClientMockRecorder) IsInitialized() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInitialized", reflect.TypeOf((*MockProjectClient)(nil).IsInitialized))
}

// Read mocks base method.
func (m *MockProjectClient) Read() (spec.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(spec.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockProjectClientMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockProjectClient)(nil).Read))
}

// Write mocks base method.
func (m *MockProjectClient) Write(projectSpec spec.Project) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", projectSpec)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockProjectClientMockRecorder) Write(projectSpec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockProjectClient)(nil).Write), projectSpec)
}

// MockConfigClient is a mock of ConfigClient interface.
type MockConfigClient struct {
	ctrl     *gomock.Controller
	recorder *MockConfigClientMockRecorder
}

// MockConfigClientMockRecorder is the mock recorder for MockConfigClient.
type MockConfigClientMockRecorder struct {
	mock *MockConfigClient
}

// NewMockConfigClient creates a new mock instance.
func NewMockConfigClient(ctrl *gomock.Controller) *MockConfigClient {
	mock := &MockConfigClient{ctrl: ctrl}
	mock.recorder = &MockConfigClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigClient) EXPECT() *MockConfigClientMockRecorder {
	return m.recorder
}

// GetFormat mocks base method.
func (m *MockConfigClient) GetFormat() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFormat")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFormat indicates an expected call of GetFormat.
func (mr *MockConfigClientMockRecorder) GetFormat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFormat", reflect.TypeOf((*MockConfigClient)(nil).GetFormat))
}

// GetUserEmailAddress mocks base method.
func (m *MockConfigClient) GetUserEmailAddress() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserEmailAddress")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserEmailAddress indicates an expected call of GetUserEmailAddress.
func (mr *MockConfigClientMockRecorder) GetUserEmailAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserEmailAddress", reflect.TypeOf((*MockConfigClient)(nil).GetUserEmailAddress))
}

// GetUserId mocks base method.
func (m *MockConfigClient) GetUserId() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserId")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserId indicates an expected call of GetUserId.
func (mr *MockConfigClientMockRecorder) GetUserId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserId", reflect.TypeOf((*MockConfigClient)(nil).GetUserId))
}

// Read mocks base method.
func (m *MockConfigClient) Read() (config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockConfigClientMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockConfigClient)(nil).Read))
}

// SetFormat mocks base method.
func (m *MockConfigClient) SetFormat(format string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFormat", format)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetFormat indicates an expected call of SetFormat.
func (mr *MockConfigClientMockRecorder) SetFormat(format interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFormat", reflect.TypeOf((*MockConfigClient)(nil).SetFormat), format)
}

// SetUserEmailAddress mocks base method.
func (m *MockConfigClient) SetUserEmailAddress(userId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserEmailAddress", userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUserEmailAddress indicates an expected call of SetUserEmailAddress.
func (mr *MockConfigClientMockRecorder) SetUserEmailAddress(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserEmailAddress", reflect.TypeOf((*MockConfigClient)(nil).SetUserEmailAddress), userId)
}

// MockStorageClient is a mock of StorageClient interface.
type MockStorageClient struct {
	ctrl     *gomock.Controller
	recorder *MockStorageClientMockRecorder
}

// MockStorageClientMockRecorder is the mock recorder for MockStorageClient.
type MockStorageClientMockRecorder struct {
	mock *MockStorageClient
}

// NewMockStorageClient creates a new mock instance.
func NewMockStorageClient(ctrl *gomock.Controller) *MockStorageClient {
	mock := &MockStorageClient{ctrl: ctrl}
	mock.recorder = &MockStorageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageClient) EXPECT() *MockStorageClientMockRecorder {
	return m.recorder
}

// ReadAsBytes mocks base method.
func (m *MockStorageClient) ReadAsBytes(url string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAsBytes", url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAsBytes indicates an expected call of ReadAsBytes.
func (mr *MockStorageClientMockRecorder) ReadAsBytes(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAsBytes", reflect.TypeOf((*MockStorageClient)(nil).ReadAsBytes), url)
}

// ReadAsString mocks base method.
func (m *MockStorageClient) ReadAsString(url string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAsString", url)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAsString indicates an expected call of ReadAsString.
func (mr *MockStorageClientMockRecorder) ReadAsString(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAsString", reflect.TypeOf((*MockStorageClient)(nil).ReadAsString), url)
}

// WriteFromBytes mocks base method.
func (m *MockStorageClient) WriteFromBytes(url string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFromBytes", url, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFromBytes indicates an expected call of WriteFromBytes.
func (mr *MockStorageClientMockRecorder) WriteFromBytes(url, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFromBytes", reflect.TypeOf((*MockStorageClient)(nil).WriteFromBytes), url, data)
}

// WriteFromString mocks base method.
func (m *MockStorageClient) WriteFromString(url, data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFromString", url, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFromString indicates an expected call of WriteFromString.
func (mr *MockStorageClientMockRecorder) WriteFromString(url, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFromString", reflect.TypeOf((*MockStorageClient)(nil).WriteFromString), url, data)
}

// MockInputClient is a mock of InputClient interface.
type MockInputClient struct {
	ctrl     *gomock.Controller
	recorder *MockInputClientMockRecorder
}

// MockInputClientMockRecorder is the mock recorder for MockInputClient.
type MockInputClientMockRecorder struct {
	mock *MockInputClient
}

// NewMockInputClient creates a new mock instance.
func NewMockInputClient(ctrl *gomock.Controller) *MockInputClient {
	mock := &MockInputClient{ctrl: ctrl}
	mock.recorder = &MockInputClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInputClient) EXPECT() *MockInputClientMockRecorder {
	return m.recorder
}

// UpdateInputReferencesAndUploadToS3 mocks base method.
func (m *MockInputClient) UpdateInputReferencesAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInputReferencesAndUploadToS3", initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInputReferencesAndUploadToS3 indicates an expected call of UpdateInputReferencesAndUploadToS3.
func (mr *MockInputClientMockRecorder) UpdateInputReferencesAndUploadToS3(initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInputReferencesAndUploadToS3", reflect.TypeOf((*MockInputClient)(nil).UpdateInputReferencesAndUploadToS3), initialProjectDirectory, tempProjectDirectory, bucketName, baseS3Key)
}
