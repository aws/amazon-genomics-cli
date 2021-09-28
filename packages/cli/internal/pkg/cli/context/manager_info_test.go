package context

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

func TestManager_Info(t *testing.T) {
	testCases := map[string]struct {
		setupMocks   func(*testing.T) mockClients
		expectedInfo Detail
		expectedErr  error
	}{
		"unstarted context": {
			expectedInfo: Detail{
				Summary: Summary{
					Name:          testContextName1,
					IsSpot:        true,
					InstanceTypes: []string{"c5"},
				},
				Status:         StatusNotStarted,
				BucketLocation: "s3://test-output-bucket/project/testProjectName/userid/bender123/context/testContextName1",
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName, Contexts: map[string]spec.Context{testContextName1: {RequestSpotInstances: true, InstanceTypes: []string{"c5"}, Engines: []spec.Engine{{Type: "wdl", Engine: "cromwell"}}}}}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.cfnMock.EXPECT().GetStackInfo("Agc-Context-testProjectName-bender123-testContextName1").
					Return(cfn.StackInfo{}, cfn.StackDoesNotExistError)
				return mockClients
			},
		},
		"started context": {
			expectedInfo: Detail{
				Summary:            Summary{Name: testContextName1},
				Status:             StatusStarted,
				BucketLocation:     "s3://test-output-bucket/project/testProjectName/userid/bender123/context/testContextName1",
				WesUrl:             testWesUrl,
				EngineLogGroupName: testLogGroupName,
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName, Contexts: map[string]spec.Context{testContextName1: {Engines: []spec.Engine{{Type: "wdl", Engine: "cromwell"}}}}}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.cfnMock.EXPECT().GetStackInfo("Agc-Context-testProjectName-bender123-testContextName1").
					Return(cfn.StackInfo{Status: types.StackStatusCreateComplete, Outputs: map[string]string{"WesUrl": testWesUrl, "EngineLogGroupName": testLogGroupName}}, nil)
				return mockClients
			},
		},
		"failed context": {
			expectedInfo: Detail{
				Summary:            Summary{Name: testContextName1},
				Status:             StatusFailed,
				BucketLocation:     "s3://test-output-bucket/project/testProjectName/userid/bender123/context/testContextName1",
				WesUrl:             testWesUrl,
				EngineLogGroupName: testLogGroupName,
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName, Contexts: map[string]spec.Context{testContextName1: {Engines: []spec.Engine{{Type: "wdl", Engine: "cromwell"}}}}}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.cfnMock.EXPECT().GetStackInfo("Agc-Context-testProjectName-bender123-testContextName1").
					Return(cfn.StackInfo{Status: types.StackStatusCreateFailed, Outputs: map[string]string{"WesUrl": testWesUrl, "EngineLogGroupName": testLogGroupName}}, nil)
				return mockClients
			},
		},
		"stopped context": {
			expectedInfo: Detail{
				Summary:            Summary{Name: testContextName1},
				Status:             StatusStopped,
				BucketLocation:     "s3://test-output-bucket/project/testProjectName/userid/bender123/context/testContextName1",
				WesUrl:             testWesUrl,
				EngineLogGroupName: testLogGroupName,
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName, Contexts: map[string]spec.Context{testContextName1: {Engines: []spec.Engine{{Type: "wdl", Engine: "cromwell"}}}}}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.cfnMock.EXPECT().GetStackInfo("Agc-Context-testProjectName-bender123-testContextName1").
					Return(cfn.StackInfo{Status: types.StackStatusDeleteComplete, Outputs: map[string]string{"WesUrl": testWesUrl, "EngineLogGroupName": testLogGroupName}}, nil)
				return mockClients
			},
		},
		"unknown context": {
			expectedErr: actionable.NewError(
				fmt.Errorf("context 'testContextName1' is not defined in Project 'testProjectName' specification"),
				"Please add the context to your project spec and deploy it or specify a different context from the command 'agc context list'",
			),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.cfnMock.EXPECT().GetStackInfo("Agc-Context-testProjectName-bender123-testContextName1").
					Return(cfn.StackInfo{Status: "unknown", Outputs: map[string]string{"EngineLogGroupName": testLogGroupName}}, nil)
				return mockClients
			},
		},
		"read error": {
			expectedErr: fmt.Errorf("some read error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{}, fmt.Errorf("some read error"))
				return mockClients
			},
		},
		"output bucket error": {
			expectedErr: fmt.Errorf("some output bucket error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return("", fmt.Errorf("some output bucket error"))
				return mockClients
			},
		},
		"stack error": {
			expectedErr: fmt.Errorf("some stack error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName, Contexts: map[string]spec.Context{testContextName1: {Engines: []spec.Engine{{Type: "wdl", Engine: "cromwell"}}}}}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.cfnMock.EXPECT().GetStackInfo("Agc-Context-testProjectName-bender123-testContextName1").
					Return(cfn.StackInfo{}, fmt.Errorf("some stack error"))
				return mockClients
			},
		},
		"context not exist error": {
			expectedErr: actionable.NewError(
				fmt.Errorf("context 'testContextName1' is not defined in Project 'testProjectName' specification"),
				"Please add the context to your project spec and deploy it or specify a different context from the command 'agc context list'",
			),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{Name: testProjectName, Contexts: map[string]spec.Context{testUnknownContextName: {Engines: []spec.Engine{{Type: "wdl", Engine: "cromwell"}}}}}, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.cfnMock.EXPECT().GetStackInfo("Agc-Context-testProjectName-bender123-testContextName1").
					Return(cfn.StackInfo{}, cfn.StackDoesNotExistError)
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
				Ssm:     mockClients.ssmMock,
				Project: mockClients.projMock,
				Config:  mockClients.configMock,
			}

			info, err := manager.Info(testContextName1)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInfo, info)
			}

		})
	}
}
