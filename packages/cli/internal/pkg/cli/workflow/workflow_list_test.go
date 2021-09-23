package workflow

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/spec"
	awsmocks "github.com/aws/amazon-genomics-cli/cli/internal/pkg/mocks/aws"
	storagemocks "github.com/aws/amazon-genomics-cli/cli/internal/pkg/mocks/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

const (
	testWorkflow1 = "Test Workflow 1"
	testWorkflow2 = "Test Workflow 2"
	testWorkflow3 = "Test Workflow 3"
)

type WorkflowListTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockProjectClient *storagemocks.MockProjectClient
	mockDdb           *awsmocks.MockDdbClient
	mockStorageClient *storagemocks.MockStorageClient
	mockConfigClient  *storagemocks.MockConfigClient

	testProjSpec spec.Project

	manager *Manager
}

func (s *WorkflowListTestSuite) BeforeTest(_, _ string) {
	s.ctrl = gomock.NewController(s.T())
	s.mockProjectClient = storagemocks.NewMockProjectClient(s.ctrl)
	s.mockDdb = awsmocks.NewMockDdbClient(s.ctrl)
	s.mockStorageClient = storagemocks.NewMockStorageClient(s.ctrl)
	s.mockConfigClient = storagemocks.NewMockConfigClient(s.ctrl)

	s.manager = &Manager{
		Project: s.mockProjectClient,
		Ddb:     s.mockDdb,
		Storage: s.mockStorageClient,
		Config:  s.mockConfigClient,
	}

	s.testProjSpec = spec.Project{
		Name: testProjectName,
		Workflows: map[string]spec.Workflow{
			testWorkflow1: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowLocalUrl,
			},
			testWorkflow2: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowS3Url,
			},
			testWorkflow3: {
				Type:      testWorkflowType,
				SourceURL: testWorkflowInvalidUrl,
			},
		},
	}
}

func (s *WorkflowListTestSuite) TestListWorkflows_EmptyList() {
	defer s.ctrl.Finish()
	s.testProjSpec.Workflows = map[string]spec.Workflow{}
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)

	summaries, err := s.manager.ListWorkflows()
	if s.Assert().NoError(err) {
		s.Assert().Empty(summaries)
	}
}

func (s *WorkflowListTestSuite) TestListWorkflows_Nominal() {
	defer s.ctrl.Finish()
	s.mockConfigClient.EXPECT().GetUserId().Return(testUserId, nil)
	s.mockProjectClient.EXPECT().Read().Return(s.testProjSpec, nil)

	actualSummaries, err := s.manager.ListWorkflows()
	expectedSummaries := map[string]Summary{
		testWorkflow1: {Name: testWorkflow1},
		testWorkflow2: {Name: testWorkflow2},
		testWorkflow3: {Name: testWorkflow3},
	}
	if s.Assert().NoError(err) {
		s.Assert().Equal(expectedSummaries, actualSummaries)
	}
}

func (s *WorkflowListTestSuite) TestListWorkflows_ReadProjectSpecFailure() {
	defer s.ctrl.Finish()
	errorMessage := "failed to read project specification"
	s.mockProjectClient.EXPECT().Read().Return(spec.Project{}, errors.New(errorMessage))

	summaries, err := s.manager.ListWorkflows()
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
		s.Assert().Empty(summaries)
	}
}

func TestWorkflowListTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowListTestSuite))
}
