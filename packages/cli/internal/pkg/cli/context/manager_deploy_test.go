package context

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestManager_Deploy(t *testing.T) {
	origVerbose := logging.Verbose
	defer func() { logging.Verbose = origVerbose }()
	logging.Verbose = true

	testCases := map[string]struct {
		setupMocks  func(*testing.T) mockClients
		expectedErr error
	}{
		"deploy success": {
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				defer close(mockClients.progressStream1)
				defer close(mockClients.progressStream2)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				clearContext := mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				mockClients.cdkMock.EXPECT().DeployApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Any()).After(clearContext).Return(mockClients.progressStream1, nil)
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
			expectedErr: fmt.Errorf("some outbut bucket error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return("", fmt.Errorf("some outbut bucket error"))
				return mockClients
			},
		},
		"context error": {
			expectedErr: fmt.Errorf("context 'testContextName1' is not defined in Project 'testProjectName' specification"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				projSpec := testValidProjectSpec
				projSpec.Contexts = nil
				mockClients.projMock.EXPECT().Read().Return(projSpec, nil)
				return mockClients
			},
		},
		"artifact bucket error": {
			expectedErr: fmt.Errorf("some artifact bucket error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return("", fmt.Errorf("some artifact bucket error"))
				return mockClients
			},
		},
		"deploy error": {
			expectedErr: fmt.Errorf("some context error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(nil)
				mockClients.cdkMock.EXPECT().DeployApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Any()).Return(nil, fmt.Errorf("some context error"))
				return mockClients
			},
		},
		"clear context error": {
			expectedErr: fmt.Errorf("failed to clear context"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.ssmMock.EXPECT().GetOutputBucket().Return(testOutputBucket, nil)
				mockClients.ssmMock.EXPECT().GetCommonParameter("installed-artifacts/s3-root-url").Return(testArtifactBucket, nil)
				mockClients.cdkMock.EXPECT().ClearContext(filepath.Join(testHomeDir, ".agc/cdk/apps/context")).Return(fmt.Errorf("failed to clear context"))
				return mockClients
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockClients := tc.setupMocks(t)
			defer mockClients.ctrl.Finish()
			manager := Manager{
				Cdk:       mockClients.cdkMock,
				Project:   mockClients.projMock,
				Ssm:       mockClients.ssmMock,
				Config:    mockClients.configMock,
				Cfn:       mockClients.cfnMock,
				baseProps: baseProps{homeDir: testHomeDir},
			}

			err := manager.Deploy(testContextName1, false)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
