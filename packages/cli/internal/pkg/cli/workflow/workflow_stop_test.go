package workflow

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	wesmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/wes"
	"github.com/aws/amazon-genomics-cli/internal/pkg/wes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type WorkflowStopTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockProjectClient *storagemocks.MockProjectClient
	mockDdb           *awsmocks.MockDdbClient
	mockConfigClient  *storagemocks.MockConfigClient
	mockCfn           *awsmocks.MockCfnClient
	mockWes           *wesmocks.MockWesClient
	testProjSpec      spec.Project
	testStackInfo     cfn.StackInfo
	wfInstance        ddb.WorkflowInstance
	manager           *Manager
}

func (s *WorkflowStopTestSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.mockProjectClient = storagemocks.NewMockProjectClient(s.ctrl)
	s.mockConfigClient = storagemocks.NewMockConfigClient(s.ctrl)
	s.mockDdb = awsmocks.NewMockDdbClient(s.ctrl)
	s.mockWes = wesmocks.NewMockWesClient(s.ctrl)
	s.mockCfn = awsmocks.NewMockCfnClient(s.ctrl)

	s.manager = &Manager{
		Project:    s.mockProjectClient,
		Ddb:        s.mockDdb,
		Config:     s.mockConfigClient,
		Cfn:        s.mockCfn,
		WesFactory: func(_ string) (wes.Interface, error) { return s.mockWes, nil },
	}

	s.testProjSpec = spec.Project{
		Name: testProjectName,
		Workflows: map[string]spec.Workflow{
			testLocalWorkflowName: {
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

	s.testStackInfo = cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWesUrl},
	}

	s.wfInstance = ddb.WorkflowInstance{
		RunId:        testRun1Id,
		WorkflowName: testLocalWorkflowName,
		ContextName:  testContext1Name,
		ProjectName:  testProjectName,
		UserId:       testUserId,
	}
}

func (s *WorkflowStopTestSuite) TestStopWorkflow_RunIdFound() {
	defer s.ctrl.Finish()

	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(context.Background(), s.testProjSpec.Name, testUserId, testRun1Id).Return(s.wfInstance, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().StopWorkflow(context.Background(), testRun1Id).Return(nil)

	s.manager.StopWorkflowInstance(testRun1Id)
	s.Assert().NoError(s.manager.err)
}

func (s *WorkflowStopTestSuite) TestStopWorkflowStop_RunIdNotFound() {
	defer s.ctrl.Finish()

	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(context.Background(), s.testProjSpec.Name, testUserId, "does-not-exist").Return(ddb.WorkflowInstance{}, fmt.Errorf("not found"))

	s.manager.StopWorkflowInstance("does-not-exist")
	s.Assert().Error(s.manager.err)
}

func (s *WorkflowStopTestSuite) TestWorkflowStop_WesUnableToStopWorkflow() {
	defer s.ctrl.Finish()

	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(context.Background(), s.testProjSpec.Name, testUserId, testRun1Id).Return(s.wfInstance, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(s.testStackInfo, nil)
	s.mockWes.EXPECT().StopWorkflow(context.Background(), testRun1Id).Return(fmt.Errorf("wf engine can't stop instance"))

	s.manager.StopWorkflowInstance(testRun1Id)
	s.Assert().Error(s.manager.err)
}

func TestWorkflowStopTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowStopTestSuite))
}
