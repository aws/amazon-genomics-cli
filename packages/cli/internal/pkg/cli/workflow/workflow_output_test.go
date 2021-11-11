package workflow

import (
	ctx "context"
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
	"github.com/rsc/wes_client"
	"github.com/stretchr/testify/suite"
)

type WorkflowOutputTestSuite struct {
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

func (s *WorkflowOutputTestSuite) BeforeTest(_, _ string) {
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

func (s *WorkflowOutputTestSuite) TestOutputByInstanceId_InstanceFound() {

	expectedOutput := map[string]interface{}{"foo": "baa"}

	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRun1Id).Return(ddb.WorkflowInstance{
		RunId:        testRun1Id,
		WorkflowName: testWorkflow1,
		ContextName:  testContext1Name,
		ProjectName:  testProjectName,
		UserId:       testUserId,
		CreatedTime:  testWorkflowSubmitTime1,
	}, nil)
	s.mockCfn.EXPECT().GetStackInfo(testContext1Stack).Return(cfn.StackInfo{
		Outputs: map[string]string{"WesUrl": testWes1Url},
	}, nil)
	s.mockWes1.EXPECT().GetRunLog(ctx.Background(), testRun1Id).Return(wes_client.RunLog{
		Outputs: map[string]interface{}{"foo": "baa"},
	}, nil)
	actualOutput, err := s.manager.OutputByInstanceId(testRun1Id)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expectedOutput, actualOutput)
	}
}

func (s *WorkflowOutputTestSuite) TestOutputByInstanceId_NoInstanceFound() {

	expectedErr := fmt.Errorf("workflow instance with id '%s' does not exist", testRun2Id)

	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockDdb.EXPECT().GetWorkflowInstanceById(ctx.Background(), testProjectName, testUserId, testRun2Id).Return(ddb.WorkflowInstance{}, fmt.Errorf("workflow instance with id '%s' does not exist", testRun2Id))
	_, err := s.manager.OutputByInstanceId(testRun2Id)
	s.Assert().EqualError(err, expectedErr.Error())

}

func TestWorkflowOutputTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowOutputTestSuite))
}
