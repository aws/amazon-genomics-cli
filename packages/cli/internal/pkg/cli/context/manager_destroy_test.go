package context

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestManager_Destroy(t *testing.T) {
	contextList := []string{testContextName1}
	origVerbose := logging.Verbose
	defer func() { logging.Verbose = origVerbose }()
	logging.Verbose = false
	testCases := map[string]struct {
		setupMocks                 func(*testing.T) mockClients
		expectedProgressResultList []ProgressResult
		contextList                []string
	}{
		"destroy success": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Outputs: []string{"some message"}, Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				defer close(mockClients.progressStream1)
				defer close(mockClients.progressStream2)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.cdkMock.EXPECT().DestroyApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Any(), testContextName1).Return(mockClients.progressStream1, nil)
				mockClients.cdkMock.EXPECT().DisplayProgressBar(fmt.Sprintf("Destroying resources for context(s) %s", contextList), []cdk.ProgressStream{mockClients.progressStream1}).Return([]cdk.Result{{Outputs: []string{"some message"}, UniqueKey: testContextName1}})
				return mockClients
			},
		},
		"read error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("some read error"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.projMock.EXPECT().Read().Return(spec.Project{}, fmt.Errorf("some read error"))
				return mockClients
			},
		},
		"destroy error": {
			contextList: contextList,
			expectedProgressResultList: []ProgressResult{
				{Err: fmt.Errorf("some context error"), Context: testContextName1},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.cdkMock.EXPECT().DestroyApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Any(), testContextName1).Return(nil, fmt.Errorf("some context error"))
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
				baseProps: baseProps{homeDir: testHomeDir},
				Config:    mockClients.configMock,
			}

			progressResultList := manager.Destroy(tc.contextList)

			if len(progressResultList) != len(tc.expectedProgressResultList) {
				assert.Equal(t, tc.expectedProgressResultList, progressResultList)
			}

			for i, progressResult := range progressResultList {
				expectedProgressResult := tc.expectedProgressResultList[i]
				assert.Equal(t, expectedProgressResult.Context, progressResult.Context)
				if expectedProgressResult.Err != nil {
					assert.Error(t, progressResult.Err, expectedProgressResult)
				} else {
					assert.NoError(t, progressResult.Err)
				}
			}
		})
	}
}
