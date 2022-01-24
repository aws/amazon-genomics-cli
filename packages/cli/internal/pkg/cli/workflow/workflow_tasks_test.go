package workflow

import (
	ctx "context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	wesmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/wes"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes"
	"github.com/golang/mock/gomock"
	"github.com/rsc/wes_client"
	"github.com/stretchr/testify/suite"
)

const (
	testRunId     = "test-run-id"
	testTaskName  = "test-task-name"
	testTaskJobId = "test-task-job-id"
	testExitCode  = "0"
)

var (
	testTaskCompositeName = fmt.Sprintf("%s|%s", testTaskName, testTaskJobId)
	testStartTime         = time.Now().Add(-time.Hour)
	testStopTime          = time.Now()
)

type GetWorkflowTasksTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockProjectClient *storagemocks.MockProjectClient
	mockDdb           *awsmocks.MockDdbClient
	mockConfigClient  *storagemocks.MockConfigClient
	mockCfn           *awsmocks.MockCfnClient
	mockWes           *wesmocks.MockWesClient

	testProjSpec spec.Project

	manager *Manager
}

func (s *GetWorkflowTasksTestSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.mockProjectClient = storagemocks.NewMockProjectClient(s.ctrl)
	s.mockDdb = awsmocks.NewMockDdbClient(s.ctrl)
	s.mockConfigClient = storagemocks.NewMockConfigClient(s.ctrl)
	s.mockCfn = awsmocks.NewMockCfnClient(s.ctrl)
	s.mockWes = wesmocks.NewMockWesClient(s.ctrl)

	wesMap := map[string]*wesmocks.MockWesClient{testWes1Url: s.mockWes}

	s.manager = &Manager{
		Project:    s.mockProjectClient,
		Ddb:        s.mockDdb,
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
		},
	}
}

func (s *GetWorkflowTasksTestSuite) TestGetWorkflowTasks_WithTask() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRunId).Return(ddb.WorkflowInstance{ContextName: testContext1Name}, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{Outputs: map[string]string{"WesUrl": testWes1Url}}, nil)
	s.mockWes.EXPECT().GetRunLog(ctx.Background(), testRunId).Return(wes_client.RunLog{
		RunId: testRunId,
		TaskLogs: []wes_client.Log{{
			Name:      testTaskCompositeName,
			StartTime: testStartTime.UTC().Format("2006-01-02T15:04:05Z"),
			EndTime:   testStopTime.UTC().Format("2006-01-02T15:04:05Z"),
			ExitCode:  testExitCode,
		}},
	}, nil)

	tasks, err := s.manager.GetWorkflowTasks(testRunId)
	if s.Assert().NoError(err) {
		s.Assert().Equal(testTaskName, tasks[0].Name)
		s.Assert().Equal(testTaskJobId, tasks[0].JobId)
		s.Assert().True(tasks[0].StartTime.Equal(testStartTime.Truncate(time.Second)))
		s.Assert().True(tasks[0].StopTime.Equal(testStopTime.Truncate(time.Second)))
		s.Assert().Equal(testExitCode, tasks[0].ExitCode)
	}
}

func (s *GetWorkflowTasksTestSuite) TestGetWorkflowTasks_ReadProjectSpecFailure() {
	defer s.ctrl.Finish()
	errorMessage := "failed to read project specification"
	s.mockProjectClient.EXPECT().Read().Return(spec.Project{}, errors.New(errorMessage))

	tasks, err := s.manager.GetWorkflowTasks(testRunId)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(tasks)
	}
}

func (s *GetWorkflowTasksTestSuite) TestGetWorkflowTasks_ReadConfigFailure() {
	defer s.ctrl.Finish()
	errorMessage := "failed to read config specification"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, errors.New(errorMessage))

	tasks, err := s.manager.GetWorkflowTasks(testRunId)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(tasks)
	}
}

func (s *GetWorkflowTasksTestSuite) TestGetWorkflowTasks_GetWorkflowInstanceFailure() {
	defer s.ctrl.Finish()
	errorMessage := "some ddb error"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRunId).Return(ddb.WorkflowInstance{}, errors.New(errorMessage))

	tasks, err := s.manager.GetWorkflowTasks(testRunId)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(tasks)
	}
}
func (s *GetWorkflowTasksTestSuite) TestGetWorkflowTasks_GetStackInfoFailure() {
	defer s.ctrl.Finish()
	errorMessage := "some stack error"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRunId).Return(ddb.WorkflowInstance{ContextName: testContext1Name}, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{}, errors.New(errorMessage))

	tasks, err := s.manager.GetWorkflowTasks(testRunId)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(tasks)
	}
}

func (s *GetWorkflowTasksTestSuite) TestGetWorkflowTasks_GetRunLogFailure() {
	defer s.ctrl.Finish()
	errorMessage := "some log error"
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRunId).Return(ddb.WorkflowInstance{ContextName: testContext1Name}, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{Outputs: map[string]string{"WesUrl": testWes1Url}}, nil)
	s.mockWes.EXPECT().GetRunLog(ctx.Background(), testRunId).Return(wes_client.RunLog{}, errors.New(errorMessage))

	tasks, err := s.manager.GetWorkflowTasks(testRunId)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(tasks)
	}
}

func (s *GetWorkflowTasksTestSuite) TestGetWorkflowTasks_WithTaskNameFailure() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRunId).Return(ddb.WorkflowInstance{ContextName: testContext1Name}, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{Outputs: map[string]string{"WesUrl": testWes1Url}}, nil)
	s.mockWes.EXPECT().GetRunLog(ctx.Background(), testRunId).Return(wes_client.RunLog{
		RunId: testRunId,
		TaskLogs: []wes_client.Log{{
			Name: testTaskName,
		}},
	}, nil)

	tasks, err := s.manager.GetWorkflowTasks(testRunId)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, "unable to parse job ID from task name 'test-task-name'")
		s.Assert().Empty(tasks)
	}
}

func TestGetWorkflowTasksTestSuite(t *testing.T) {
	suite.Run(t, new(GetWorkflowTasksTestSuite))
}
