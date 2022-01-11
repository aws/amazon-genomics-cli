package workflow

import (
	"errors"
	"testing"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	storagemocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/storage"
	wesmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/wes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type WorkflowDescribeTestSuite struct {
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

func (s *WorkflowDescribeTestSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.mockProjectClient = storagemocks.NewMockProjectClient(s.ctrl)
	s.mockDdb = awsmocks.NewMockDdbClient(s.ctrl)
	s.mockStorageClient = storagemocks.NewMockStorageClient(s.ctrl)
	s.mockConfigClient = storagemocks.NewMockConfigClient(s.ctrl)
	s.mockCfn = awsmocks.NewMockCfnClient(s.ctrl)
	s.mockWes1 = wesmocks.NewMockWesClient(s.ctrl)
	s.mockWes2 = wesmocks.NewMockWesClient(s.ctrl)

	s.manager = &Manager{
		Project: s.mockProjectClient,
		Ddb:     s.mockDdb,
		Storage: s.mockStorageClient,
		Config:  s.mockConfigClient,
		Cfn:     s.mockCfn,
	}

	s.testProjSpec = spec.Project{
		Name: testProjectName,
		Workflows: map[string]spec.Workflow{
			testWorkflow1: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowLocalUrl,
			},
		},
	}
}

func (s *WorkflowDescribeTestSuite) TestDescribeWorkflow_Nominal() {
	defer s.ctrl.Finish()
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)

	actualDetails, err := s.manager.DescribeWorkflow(testWorkflow1)
	if s.Assert().NoError(err) {
		s.Assert().Equal(testWorkflow1, actualDetails.Name)
		s.Assert().Equal(testWorkflowLocalUrl, actualDetails.Source)
		s.Assert().Equal(testWorkflowTypeLang, actualDetails.TypeLanguage)
		s.Assert().Equal(testWorkflowTypeVer, actualDetails.TypeVersion)
	}
}

func (s *WorkflowDescribeTestSuite) TestDescribeWorkflow_ReadProjectSpecFailure() {
	defer s.ctrl.Finish()
	errorMessage := "failed to read project specification"
	s.mockProjectClient.EXPECT().Read().Return(spec.Project{}, errors.New(errorMessage))

	actualDetails, err := s.manager.DescribeWorkflow(testWorkflow1)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(actualDetails)
	}
}

func TestWorkflowDescribeTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowDescribeTestSuite))
}
