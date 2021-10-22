// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/mocks/aws/interfaces.go

// Package awsmocks is a generated GoMock package.
package awsmocks

import (
	context "context"
	reflect "reflect"
	regexp "regexp"

	batch "github.com/aws/amazon-genomics-cli/internal/pkg/aws/batch"
	cdk "github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	cfn "github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	cwl "github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
	ddb "github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	ecr "github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	types "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	gomock "github.com/golang/mock/gomock"
)

// MockCdkClient is a mock of CdkClient interface.
type MockCdkClient struct {
	ctrl     *gomock.Controller
	recorder *MockCdkClientMockRecorder
}

// MockCdkClientMockRecorder is the mock recorder for MockCdkClient.
type MockCdkClientMockRecorder struct {
	mock *MockCdkClient
}

// NewMockCdkClient creates a new mock instance.
func NewMockCdkClient(ctrl *gomock.Controller) *MockCdkClient {
	mock := &MockCdkClient{ctrl: ctrl}
	mock.recorder = &MockCdkClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCdkClient) EXPECT() *MockCdkClientMockRecorder {
	return m.recorder
}

// ClearContext mocks base method.
func (m *MockCdkClient) ClearContext(appDir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearContext", appDir)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearContext indicates an expected call of ClearContext.
func (mr *MockCdkClientMockRecorder) ClearContext(appDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearContext", reflect.TypeOf((*MockCdkClient)(nil).ClearContext), appDir)
}

// DeployApp mocks base method.
func (m *MockCdkClient) DeployApp(appDir string, context []string) (cdk.ProgressStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployApp", appDir, context)
	ret0, _ := ret[0].(cdk.ProgressStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeployApp indicates an expected call of DeployApp.
func (mr *MockCdkClientMockRecorder) DeployApp(appDir, context interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployApp", reflect.TypeOf((*MockCdkClient)(nil).DeployApp), appDir, context)
}

// DestroyApp mocks base method.
func (m *MockCdkClient) DestroyApp(appDir string, context []string) (cdk.ProgressStream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyApp", appDir, context)
	ret0, _ := ret[0].(cdk.ProgressStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DestroyApp indicates an expected call of DestroyApp.
func (mr *MockCdkClientMockRecorder) DestroyApp(appDir, context interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyApp", reflect.TypeOf((*MockCdkClient)(nil).DestroyApp), appDir, context)
}

// MockS3Client is a mock of S3Client interface.
type MockS3Client struct {
	ctrl     *gomock.Controller
	recorder *MockS3ClientMockRecorder
}

// MockS3ClientMockRecorder is the mock recorder for MockS3Client.
type MockS3ClientMockRecorder struct {
	mock *MockS3Client
}

// NewMockS3Client creates a new mock instance.
func NewMockS3Client(ctrl *gomock.Controller) *MockS3Client {
	mock := &MockS3Client{ctrl: ctrl}
	mock.recorder = &MockS3ClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockS3Client) EXPECT() *MockS3ClientMockRecorder {
	return m.recorder
}

// BucketExists mocks base method.
func (m *MockS3Client) BucketExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BucketExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BucketExists indicates an expected call of BucketExists.
func (mr *MockS3ClientMockRecorder) BucketExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BucketExists", reflect.TypeOf((*MockS3Client)(nil).BucketExists), arg0)
}

// SyncFile mocks base method.
func (m *MockS3Client) SyncFile(bucketName, key, filePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncFile", bucketName, key, filePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncFile indicates an expected call of SyncFile.
func (mr *MockS3ClientMockRecorder) SyncFile(bucketName, key, filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncFile", reflect.TypeOf((*MockS3Client)(nil).SyncFile), bucketName, key, filePath)
}

// UploadFile mocks base method.
func (m *MockS3Client) UploadFile(bucketName, key, filePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", bucketName, key, filePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockS3ClientMockRecorder) UploadFile(bucketName, key, filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockS3Client)(nil).UploadFile), bucketName, key, filePath)
}

// MockStsClient is a mock of StsClient interface.
type MockStsClient struct {
	ctrl     *gomock.Controller
	recorder *MockStsClientMockRecorder
}

// MockStsClientMockRecorder is the mock recorder for MockStsClient.
type MockStsClientMockRecorder struct {
	mock *MockStsClient
}

// NewMockStsClient creates a new mock instance.
func NewMockStsClient(ctrl *gomock.Controller) *MockStsClient {
	mock := &MockStsClient{ctrl: ctrl}
	mock.recorder = &MockStsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStsClient) EXPECT() *MockStsClientMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockStsClient) GetAccount() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockStsClientMockRecorder) GetAccount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockStsClient)(nil).GetAccount))
}

// MockSsmClient is a mock of SsmClient interface.
type MockSsmClient struct {
	ctrl     *gomock.Controller
	recorder *MockSsmClientMockRecorder
}

// MockSsmClientMockRecorder is the mock recorder for MockSsmClient.
type MockSsmClientMockRecorder struct {
	mock *MockSsmClient
}

// NewMockSsmClient creates a new mock instance.
func NewMockSsmClient(ctrl *gomock.Controller) *MockSsmClient {
	mock := &MockSsmClient{ctrl: ctrl}
	mock.recorder = &MockSsmClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSsmClient) EXPECT() *MockSsmClientMockRecorder {
	return m.recorder
}

// GetCommonParameter mocks base method.
func (m *MockSsmClient) GetCommonParameter(parameterSuffix string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommonParameter", parameterSuffix)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommonParameter indicates an expected call of GetCommonParameter.
func (mr *MockSsmClientMockRecorder) GetCommonParameter(parameterSuffix interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommonParameter", reflect.TypeOf((*MockSsmClient)(nil).GetCommonParameter), parameterSuffix)
}

// GetOutputBucket mocks base method.
func (m *MockSsmClient) GetOutputBucket() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutputBucket")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOutputBucket indicates an expected call of GetOutputBucket.
func (mr *MockSsmClientMockRecorder) GetOutputBucket() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutputBucket", reflect.TypeOf((*MockSsmClient)(nil).GetOutputBucket))
}

// MockCfnClient is a mock of CfnClient interface.
type MockCfnClient struct {
	ctrl     *gomock.Controller
	recorder *MockCfnClientMockRecorder
}

// MockCfnClientMockRecorder is the mock recorder for MockCfnClient.
type MockCfnClientMockRecorder struct {
	mock *MockCfnClient
}

// NewMockCfnClient creates a new mock instance.
func NewMockCfnClient(ctrl *gomock.Controller) *MockCfnClient {
	mock := &MockCfnClient{ctrl: ctrl}
	mock.recorder = &MockCfnClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCfnClient) EXPECT() *MockCfnClientMockRecorder {
	return m.recorder
}

// DeleteStack mocks base method.
func (m *MockCfnClient) DeleteStack(stackId string) (chan cfn.DeletionResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStack", stackId)
	ret0, _ := ret[0].(chan cfn.DeletionResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteStack indicates an expected call of DeleteStack.
func (mr *MockCfnClientMockRecorder) DeleteStack(stackId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStack", reflect.TypeOf((*MockCfnClient)(nil).DeleteStack), stackId)
}

// GetStackInfo mocks base method.
func (m *MockCfnClient) GetStackInfo(stackName string) (cfn.StackInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStackInfo", stackName)
	ret0, _ := ret[0].(cfn.StackInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStackInfo indicates an expected call of GetStackInfo.
func (mr *MockCfnClientMockRecorder) GetStackInfo(stackName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStackInfo", reflect.TypeOf((*MockCfnClient)(nil).GetStackInfo), stackName)
}

// GetStackOutputs mocks base method.
func (m *MockCfnClient) GetStackOutputs(stackName string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStackOutputs", stackName)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStackOutputs indicates an expected call of GetStackOutputs.
func (mr *MockCfnClientMockRecorder) GetStackOutputs(stackName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStackOutputs", reflect.TypeOf((*MockCfnClient)(nil).GetStackOutputs), stackName)
}

// GetStackStatus mocks base method.
func (m *MockCfnClient) GetStackStatus(stackName string) (types.StackStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStackStatus", stackName)
	ret0, _ := ret[0].(types.StackStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStackStatus indicates an expected call of GetStackStatus.
func (mr *MockCfnClientMockRecorder) GetStackStatus(stackName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStackStatus", reflect.TypeOf((*MockCfnClient)(nil).GetStackStatus), stackName)
}

// GetStackTags mocks base method.
func (m *MockCfnClient) GetStackTags(stackName string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStackTags", stackName)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStackTags indicates an expected call of GetStackTags.
func (mr *MockCfnClientMockRecorder) GetStackTags(stackName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStackTags", reflect.TypeOf((*MockCfnClient)(nil).GetStackTags), stackName)
}

// ListStacks mocks base method.
func (m *MockCfnClient) ListStacks(regexNameFilter *regexp.Regexp, statusFilter []types.StackStatus) ([]cfn.Stack, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStacks", regexNameFilter, statusFilter)
	ret0, _ := ret[0].([]cfn.Stack)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStacks indicates an expected call of ListStacks.
func (mr *MockCfnClientMockRecorder) ListStacks(regexNameFilter, statusFilter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStacks", reflect.TypeOf((*MockCfnClient)(nil).ListStacks), regexNameFilter, statusFilter)
}

// MockBatchClient is a mock of BatchClient interface.
type MockBatchClient struct {
	ctrl     *gomock.Controller
	recorder *MockBatchClientMockRecorder
}

// MockBatchClientMockRecorder is the mock recorder for MockBatchClient.
type MockBatchClientMockRecorder struct {
	mock *MockBatchClient
}

// NewMockBatchClient creates a new mock instance.
func NewMockBatchClient(ctrl *gomock.Controller) *MockBatchClient {
	mock := &MockBatchClient{ctrl: ctrl}
	mock.recorder = &MockBatchClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBatchClient) EXPECT() *MockBatchClientMockRecorder {
	return m.recorder
}

// GetJobs mocks base method.
func (m *MockBatchClient) GetJobs(jobIds []string) ([]batch.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobs", jobIds)
	ret0, _ := ret[0].([]batch.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobs indicates an expected call of GetJobs.
func (mr *MockBatchClientMockRecorder) GetJobs(jobIds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobs", reflect.TypeOf((*MockBatchClient)(nil).GetJobs), jobIds)
}

// MockCwlClient is a mock of CwlClient interface.
type MockCwlClient struct {
	ctrl     *gomock.Controller
	recorder *MockCwlClientMockRecorder
}

// MockCwlClientMockRecorder is the mock recorder for MockCwlClient.
type MockCwlClientMockRecorder struct {
	mock *MockCwlClient
}

// NewMockCwlClient creates a new mock instance.
func NewMockCwlClient(ctrl *gomock.Controller) *MockCwlClient {
	mock := &MockCwlClient{ctrl: ctrl}
	mock.recorder = &MockCwlClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCwlClient) EXPECT() *MockCwlClientMockRecorder {
	return m.recorder
}

// GetLogsPaginated mocks base method.
func (m *MockCwlClient) GetLogsPaginated(input cwl.GetLogsInput) cwl.LogPaginator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogsPaginated", input)
	ret0, _ := ret[0].(cwl.LogPaginator)
	return ret0
}

// GetLogsPaginated indicates an expected call of GetLogsPaginated.
func (mr *MockCwlClientMockRecorder) GetLogsPaginated(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogsPaginated", reflect.TypeOf((*MockCwlClient)(nil).GetLogsPaginated), input)
}

// StreamLogs mocks base method.
func (m *MockCwlClient) StreamLogs(ctx context.Context, logGroupName string, streams ...string) <-chan cwl.StreamEvent {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, logGroupName}
	for _, a := range streams {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StreamLogs", varargs...)
	ret0, _ := ret[0].(<-chan cwl.StreamEvent)
	return ret0
}

// StreamLogs indicates an expected call of StreamLogs.
func (mr *MockCwlClientMockRecorder) StreamLogs(ctx, logGroupName interface{}, streams ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, logGroupName}, streams...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamLogs", reflect.TypeOf((*MockCwlClient)(nil).StreamLogs), varargs...)
}

// MockCwlLogPaginator is a mock of CwlLogPaginator interface.
type MockCwlLogPaginator struct {
	ctrl     *gomock.Controller
	recorder *MockCwlLogPaginatorMockRecorder
}

// MockCwlLogPaginatorMockRecorder is the mock recorder for MockCwlLogPaginator.
type MockCwlLogPaginatorMockRecorder struct {
	mock *MockCwlLogPaginator
}

// NewMockCwlLogPaginator creates a new mock instance.
func NewMockCwlLogPaginator(ctrl *gomock.Controller) *MockCwlLogPaginator {
	mock := &MockCwlLogPaginator{ctrl: ctrl}
	mock.recorder = &MockCwlLogPaginatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCwlLogPaginator) EXPECT() *MockCwlLogPaginatorMockRecorder {
	return m.recorder
}

// HasMoreLogs mocks base method.
func (m *MockCwlLogPaginator) HasMoreLogs() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasMoreLogs")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasMoreLogs indicates an expected call of HasMoreLogs.
func (mr *MockCwlLogPaginatorMockRecorder) HasMoreLogs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasMoreLogs", reflect.TypeOf((*MockCwlLogPaginator)(nil).HasMoreLogs))
}

// NextLogs mocks base method.
func (m *MockCwlLogPaginator) NextLogs() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextLogs")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextLogs indicates an expected call of NextLogs.
func (mr *MockCwlLogPaginatorMockRecorder) NextLogs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextLogs", reflect.TypeOf((*MockCwlLogPaginator)(nil).NextLogs))
}

// MockDdbClient is a mock of DdbClient interface.
type MockDdbClient struct {
	ctrl     *gomock.Controller
	recorder *MockDdbClientMockRecorder
}

// MockDdbClientMockRecorder is the mock recorder for MockDdbClient.
type MockDdbClientMockRecorder struct {
	mock *MockDdbClient
}

// NewMockDdbClient creates a new mock instance.
func NewMockDdbClient(ctrl *gomock.Controller) *MockDdbClient {
	mock := &MockDdbClient{ctrl: ctrl}
	mock.recorder = &MockDdbClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDdbClient) EXPECT() *MockDdbClientMockRecorder {
	return m.recorder
}

// GetWorkflowInstanceById mocks base method.
func (m *MockDdbClient) GetWorkflowInstanceById(ctx context.Context, project, user, runId string) (ddb.WorkflowInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkflowInstanceById", ctx, project, user, runId)
	ret0, _ := ret[0].(ddb.WorkflowInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkflowInstanceById indicates an expected call of GetWorkflowInstanceById.
func (mr *MockDdbClientMockRecorder) GetWorkflowInstanceById(ctx, project, user, runId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkflowInstanceById", reflect.TypeOf((*MockDdbClient)(nil).GetWorkflowInstanceById), ctx, project, user, runId)
}

// ListWorkflowInstances mocks base method.
func (m *MockDdbClient) ListWorkflowInstances(ctx context.Context, project, user string, limit int) ([]ddb.WorkflowInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWorkflowInstances", ctx, project, user, limit)
	ret0, _ := ret[0].([]ddb.WorkflowInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWorkflowInstances indicates an expected call of ListWorkflowInstances.
func (mr *MockDdbClientMockRecorder) ListWorkflowInstances(ctx, project, user, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWorkflowInstances", reflect.TypeOf((*MockDdbClient)(nil).ListWorkflowInstances), ctx, project, user, limit)
}

// ListWorkflowInstancesByContext mocks base method.
func (m *MockDdbClient) ListWorkflowInstancesByContext(ctx context.Context, project, user, contextName string, limit int) ([]ddb.WorkflowInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWorkflowInstancesByContext", ctx, project, user, contextName, limit)
	ret0, _ := ret[0].([]ddb.WorkflowInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWorkflowInstancesByContext indicates an expected call of ListWorkflowInstancesByContext.
func (mr *MockDdbClientMockRecorder) ListWorkflowInstancesByContext(ctx, project, user, contextName, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWorkflowInstancesByContext", reflect.TypeOf((*MockDdbClient)(nil).ListWorkflowInstancesByContext), ctx, project, user, contextName, limit)
}

// ListWorkflowInstancesByName mocks base method.
func (m *MockDdbClient) ListWorkflowInstancesByName(ctx context.Context, project, user, workflowName string, limit int) ([]ddb.WorkflowInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWorkflowInstancesByName", ctx, project, user, workflowName, limit)
	ret0, _ := ret[0].([]ddb.WorkflowInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWorkflowInstancesByName indicates an expected call of ListWorkflowInstancesByName.
func (mr *MockDdbClientMockRecorder) ListWorkflowInstancesByName(ctx, project, user, workflowName, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWorkflowInstancesByName", reflect.TypeOf((*MockDdbClient)(nil).ListWorkflowInstancesByName), ctx, project, user, workflowName, limit)
}

// ListWorkflows mocks base method.
func (m *MockDdbClient) ListWorkflows(ctx context.Context, project, user string) ([]ddb.WorkflowSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWorkflows", ctx, project, user)
	ret0, _ := ret[0].([]ddb.WorkflowSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWorkflows indicates an expected call of ListWorkflows.
func (mr *MockDdbClientMockRecorder) ListWorkflows(ctx, project, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWorkflows", reflect.TypeOf((*MockDdbClient)(nil).ListWorkflows), ctx, project, user)
}

// WriteWorkflowInstance mocks base method.
func (m *MockDdbClient) WriteWorkflowInstance(ctx context.Context, instance ddb.WorkflowInstance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteWorkflowInstance", ctx, instance)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteWorkflowInstance indicates an expected call of WriteWorkflowInstance.
func (mr *MockDdbClientMockRecorder) WriteWorkflowInstance(ctx, instance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteWorkflowInstance", reflect.TypeOf((*MockDdbClient)(nil).WriteWorkflowInstance), ctx, instance)
}

// MockEcrClient is a mock of EcrClient interface.
type MockEcrClient struct {
	ctrl     *gomock.Controller
	recorder *MockEcrClientMockRecorder
}

// MockEcrClientMockRecorder is the mock recorder for MockEcrClient.
type MockEcrClientMockRecorder struct {
	mock *MockEcrClient
}

// NewMockEcrClient creates a new mock instance.
func NewMockEcrClient(ctrl *gomock.Controller) *MockEcrClient {
	mock := &MockEcrClient{ctrl: ctrl}
	mock.recorder = &MockEcrClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEcrClient) EXPECT() *MockEcrClientMockRecorder {
	return m.recorder
}

// ImageListable mocks base method.
func (m *MockEcrClient) ImageListable(arg0, arg1, arg2, arg3 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageListable", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageListable indicates an expected call of ImageListable.
func (mr *MockEcrClientMockRecorder) ImageListable(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageListable", reflect.TypeOf((*MockEcrClient)(nil).ImageListable), arg0, arg1, arg2, arg3)
}

// VerifyImageExists mocks base method.
func (m *MockEcrClient) VerifyImageExists(reference ecr.ImageReference) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyImageExists", reference)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyImageExists indicates an expected call of VerifyImageExists.
func (mr *MockEcrClientMockRecorder) VerifyImageExists(reference interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyImageExists", reflect.TypeOf((*MockEcrClient)(nil).VerifyImageExists), reference)
}
