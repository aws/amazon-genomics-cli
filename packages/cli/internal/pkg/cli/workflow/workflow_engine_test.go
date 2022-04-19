package workflow

import (
	ctx "context"
	"testing"

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

type WorkflowEngineTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockProjectClient *storagemocks.MockProjectClient
	mockDdbClient     *awsmocks.MockDdbClient
	mockStorageClient *storagemocks.MockStorageClient
	mockConfigClient  *storagemocks.MockConfigClient
	mockCfnClient     *awsmocks.MockCfnClient
	mockWesClient     *wesmocks.MockWesClient

	testProjSpec spec.Project
	manager      *Manager
}

func (s *WorkflowEngineTestSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.mockProjectClient = storagemocks.NewMockProjectClient(s.ctrl)
	s.mockDdbClient = awsmocks.NewMockDdbClient(s.ctrl)
	s.mockStorageClient = storagemocks.NewMockStorageClient(s.ctrl)
	s.mockConfigClient = storagemocks.NewMockConfigClient(s.ctrl)
	s.mockCfnClient = awsmocks.NewMockCfnClient(s.ctrl)
	s.mockWesClient = wesmocks.NewMockWesClient(s.ctrl)
	s.manager = &Manager{
		Project:     s.mockProjectClient,
		Config:      s.mockConfigClient,
		S3:          nil,
		Ssm:         nil,
		Cfn:         s.mockCfnClient,
		Ddb:         s.mockDdbClient,
		Storage:     s.mockStorageClient,
		InputClient: nil,
		WesFactory:  func(url string) (wes.Interface, error) { return s.mockWesClient, nil },
	}
	s.testProjSpec = spec.Project{
		Name:      testProjectName,
		Workflows: nil,
		Data:      nil,
		Contexts: map[string]spec.Context{
			testContext1Name: {Engines: []spec.Engine{{Type: "wdl", Engine: "miniwdl"}}},
		},
	}
}

func (s *WorkflowEngineTestSuite) Test_GetEngineLogByRunId_Success() {
	expectedOutput := "cloudWatchStreamName"

	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdbClient.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRun1Id).Return(ddb.WorkflowInstance{
		RunId:        testRun1Id,
		WorkflowName: testWorkflow1,
		ContextName:  testContext1Name,
		ProjectName:  testProjectName,
		UserId:       testUserId,
		CreatedTime:  testWorkflowSubmitTime1,
	}, nil)
	s.mockCfnClient.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWesClient.EXPECT().GetRunLog(ctx.Background(), testRun1Id).Return(wes_client.RunLog{
		RunId: testRun1Id,
		RunLog: wes_client.Log{
			Stdout: "cloudWatchStreamName",
		},
	}, nil)

	engineLog, err := s.manager.GetEngineLogByRunId(testRun1Id)
	actualOutput := engineLog.StdOut
	if s.Assert().NoError(err) {
		s.Assert().Equal(expectedOutput, actualOutput)
	}
}

func TestWorkflowEngineTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowEngineTestSuite))
}
