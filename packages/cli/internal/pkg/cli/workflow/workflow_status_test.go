package workflow

import (
	ctx "context"
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	wesmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/wes"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

const (
	testWorkflowInstancesLimit = 34
	testWorkflowSubmitTime1    = "Test Workflow Submit Time 1"
	testWorkflowSubmitTime2    = "Test Workflow Submit Time 2"
	testRunStatus1             = "TestStatus1"
	testRunStatus2             = "TestStatus2"
	testWes1Url                = "https://TestWes1Url.com/prod"
	testWes2Url                = "https://TestWes2Url.com/prod"
	testRunStatusUnknown       = "UNKNOWN"
)

var workflowInstance1 = ddb.WorkflowInstance{
	RunId:        testRun1Id,
	WorkflowName: testWorkflow1,
	ContextName:  testContext1Name,
	ProjectName:  testProjectName,
	UserId:       testUserId,
	CreatedTime:  testWorkflowSubmitTime1,
}

var workflowInstance2 = ddb.WorkflowInstance{
	RunId:        testRun2Id,
	WorkflowName: testWorkflow1,
	ContextName:  testContext1Name,
	ProjectName:  testProjectName,
	UserId:       testUserId,
	CreatedTime:  testWorkflowSubmitTime2,
}

var instanceSummary1 = InstanceSummary{
	Id:           testRun1Id,
	WorkflowName: testWorkflow1,
	ContextName:  testContext1Name,
	SubmitTime:   testWorkflowSubmitTime1,
	InProject:    true,
	State:        testRunStatus1,
}

var instanceSummary2 = InstanceSummary{
	Id:           testRun2Id,
	WorkflowName: testWorkflow1,
	ContextName:  testContext1Name,
	SubmitTime:   testWorkflowSubmitTime2,
	InProject:    true,
	State:        testRunStatus2,
}

type WorkflowStatusTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockProjectClient *storagemocks.MockProjectClient
	mockDdb           *awsmocks.MockDdbClient
	mockStorageClient *storagemocks.MockStorageClient
	mockConfigClient  *storagemocks.MockConfigClient
	mockCfn           *awsmocks.MockCfnClient
	mockWes1          *wesmocks.MockWesClient
	mockWes2          *wesmocks.MockWesClient

	testProjSpec spec.Project

	manager *Manager
}

func (s *WorkflowStatusTestSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.mockProjectClient = storagemocks.NewMockProjectClient(s.ctrl)
	s.mockDdb = awsmocks.NewMockDdbClient(s.ctrl)
	s.mockStorageClient = storagemocks.NewMockStorageClient(s.ctrl)
	s.mockConfigClient = storagemocks.NewMockConfigClient(s.ctrl)
	s.mockCfn = awsmocks.NewMockCfnClient(s.ctrl)
	s.mockWes1 = wesmocks.NewMockWesClient(s.ctrl)
	s.mockWes2 = wesmocks.NewMockWesClient(s.ctrl)

	wesMap := map[string]*wesmocks.MockWesClient{
		testWes1Url: s.mockWes1,
		testWes2Url: s.mockWes2,
	}

	s.manager = &Manager{
		Project:    s.mockProjectClient,
		Ddb:        s.mockDdb,
		Storage:    s.mockStorageClient,
		Config:     s.mockConfigClient,
		Cfn:        s.mockCfn,
		WesFactory: func(url string) (wes.Interface, error) { return wesMap[url], nil },
	}

	s.testProjSpec = spec.Project{
		Name: testProjectName,
		Workflows: map[string]spec.Workflow{
			testWorkflow1: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowLocalUrl,
			},
		},
		Contexts: map[string]spec.Context{
			testContext1Name: {
				Engines: []spec.Engine{
					{
						Type:   "wdl",
						Engine: "cromwell",
					},
				},
			},
			testContext2Name: {
				Engines: []spec.Engine{
					{
						Type:   "wdl",
						Engine: "cromwell",
					},
				},
			},
		},
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflowAll_NoInstances() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(nil, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		s.Assert().Empty(actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflowAll_InstancesSameContext() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance2,
		workflowInstance1,
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun2Id).Return(testRunStatus2, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary2,
			instanceSummary1,
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflowAll_InstancesDifferentContexts() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
		{
			RunId:        testRun2Id,
			WorkflowName: testWorkflow1,
			ContextName:  testContext2Name,
			ProjectName:  testProjectName,
			UserId:       testUserId,
			CreatedTime:  testWorkflowSubmitTime2,
		},
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext2Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext2Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes2Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)
	s.mockWes2.EXPECT().GetRunStatus(context.Background(), testRun2Id).Return(testRunStatus2, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary1,
			{
				Id:           testRun2Id,
				WorkflowName: testWorkflow1,
				ContextName:  testContext2Name,
				SubmitTime:   testWorkflowSubmitTime2,
				InProject:    true,
				State:        testRunStatus2,
			},
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflowAll_WorkflowNotInProject() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
		{
			RunId:        testRun2Id,
			WorkflowName: testWorkflow2,
			ContextName:  testContext1Name,
			ProjectName:  testProjectName,
			UserId:       testUserId,
			CreatedTime:  testWorkflowSubmitTime2,
		},
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun2Id).Return(testRunStatus2, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary1,
			{
				Id:           testRun2Id,
				WorkflowName: testWorkflow2,
				ContextName:  testContext1Name,
				SubmitTime:   testWorkflowSubmitTime2,
				InProject:    false,
				State:        testRunStatus2,
			},
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_SomeUnknownInstance() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
		workflowInstance2,
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatusUnknown, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun2Id).Return(testRunStatus2, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary2,
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_AllUnknownInstance() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatusUnknown, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		s.Assert().Empty(actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_ErrorStatusContexts() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
		{
			RunId:        testRun2Id,
			WorkflowName: testWorkflow1,
			ContextName:  testContext2Name,
			ProjectName:  testProjectName,
			UserId:       testUserId,
			CreatedTime:  testWorkflowSubmitTime2,
		},
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext2Stack).Return(types.StackStatus(""), cfn.StackDoesNotExistError)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary1,
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_NonActiveContexts() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
		{
			RunId:        testRun2Id,
			WorkflowName: testWorkflow1,
			ContextName:  testContext2Name,
			ProjectName:  testProjectName,
			UserId:       testUserId,
			CreatedTime:  testWorkflowSubmitTime2,
		},
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext2Stack).Return(types.StackStatusDeleteInProgress, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)

	actualStatuses, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary1,
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_ReadProjectSpecFailure() {
	defer s.ctrl.Finish()
	errorMessage := "failed to read project specification"
	s.mockProjectClient.EXPECT().Read().Return(spec.Project{}, errors.New(errorMessage))

	actualDetails, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(actualDetails)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_ListInstancesFailure() {
	defer s.ctrl.Finish()
	errorMessage := "failed to list workflow instances"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(nil, errors.New(errorMessage))

	actualDetails, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(actualDetails)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_CfnFailed1() {
	defer s.ctrl.Finish()
	errorMessage := "cannot call CFN stack info"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{}, errors.New(errorMessage))

	actualDetails, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(actualDetails)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_CfnFailed2() {
	defer s.ctrl.Finish()
	errorMessage := "cannot call CFN stack status"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatus(""), errors.New(errorMessage))

	actualDetails, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(actualDetails)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflow_WesFailed() {
	defer s.ctrl.Finish()
	errorMessage := "cannot call WES"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance1,
	}
	s.mockDdb.EXPECT().ListWorkflowInstances(ctx.Background(), testProjectName, testUserId, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return("", errors.New(errorMessage))

	actualDetails, err := s.manager.StatusWorkflowAll(testWorkflowInstancesLimit)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(actualDetails)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflowByContext_Nominal() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance2,
		workflowInstance1,
	}
	s.mockDdb.EXPECT().ListWorkflowInstancesByContext(ctx.Background(), testProjectName, testUserId, testContext1Name, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun2Id).Return(testRunStatus2, nil)

	actualStatuses, err := s.manager.StatusWorkflowByContext(testContext1Name, testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary2,
			instanceSummary1,
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflowByName_Nominal() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	instances := []ddb.WorkflowInstance{
		workflowInstance2,
		workflowInstance1,
	}
	s.mockDdb.EXPECT().ListWorkflowInstancesByName(ctx.Background(), testProjectName, testUserId, testWorkflow1, testWorkflowInstancesLimit).Return(instances, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun2Id).Return(testRunStatus2, nil)

	actualStatuses, err := s.manager.StatusWorkflowByName(testWorkflow1, testWorkflowInstancesLimit)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary2,
			instanceSummary1,
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func (s *WorkflowStatusTestSuite) TestStatusWorkflowByInstanceId_Nominal() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRun1Id).Return(workflowInstance1, nil)
	s.mockCfn.EXPECT().GetStackStatus(testContext1Stack).Return(types.StackStatusCreateComplete, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunStatus(context.Background(), testRun1Id).Return(testRunStatus1, nil)

	actualStatuses, err := s.manager.StatusWorkflowByInstanceId(testRun1Id)
	if s.Assert().NoError(err) {
		expectedStatuses := []InstanceSummary{
			instanceSummary1,
		}
		s.Assert().Equal(expectedStatuses, actualStatuses)
	}
}

func TestWorkflowStatusTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowStatusTestSuite))
}
