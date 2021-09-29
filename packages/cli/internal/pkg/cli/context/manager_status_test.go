package context

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/awsresources"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

func TestManager_Status(t *testing.T) {
	testCases := map[string]struct {
		setupMocks       func(*testing.T) mockClients
		expectedContexts []Instance
		expectedErr      error
	}{
		"active and inactive contexts": {
			expectedContexts: []Instance{
				{
					ContextName:            testContextName1,
					ContextStatus:          "STARTED",
					ContextReason:          "some reason",
					IsDefinedInProjectFile: true,
				},
				{
					ContextName:            testContextName2,
					ContextStatus:          "STARTED",
					ContextReason:          "other reason",
					IsDefinedInProjectFile: true,
				},
				{
					ContextName:            "testContextName45",
					ContextStatus:          "STOPPED",
					ContextReason:          "some-reason",
					IsDefinedInProjectFile: false,
				}},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				stackNamePattern := awsresources.RenderContextStackNameRegexp(testProjectName, testUserId)
				mockClients.cfnMock.EXPECT().ListStacks(regexp.MustCompile(stackNamePattern), cfn.ActiveStacksFilter).
					Return([]cfn.Stack{{
						Status:       types.StackStatusCreateComplete,
						Name:         "Agc-Context-testProjectName-bender123-testContextName1",
						StatusReason: "some reason",
					}, {
						Status:       types.StackStatusCreateComplete,
						Name:         "Agc-Context-testProjectName-bender123-testContextName2",
						StatusReason: "other reason",
					}, {
						Status:       types.StackStatusDeleteInProgress,
						Name:         "Agc-Context-testProjectName-bender123-testContextName45",
						StatusReason: "some-reason",
					}}, nil)
				return mockClients
			},
		},
		"cfn error": {
			expectedErr: fmt.Errorf("some cfn error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName}, nil)
				stackNamePattern := awsresources.RenderContextStackNameRegexp(testProjectName, testUserId)
				mockClients.cfnMock.EXPECT().ListStacks(regexp.MustCompile(stackNamePattern), cfn.ActiveStacksFilter).
					Return([]cfn.Stack{}, fmt.Errorf("some cfn error"))
				return mockClients
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockClients := tc.setupMocks(t)
			defer mockClients.ctrl.Finish()
			manager := Manager{
				Cfn:     mockClients.cfnMock,
				Project: mockClients.projMock,
				Config:  mockClients.configMock,
			}

			contexts, err := manager.StatusList()

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedContexts, contexts)
			}

		})
	}
}
