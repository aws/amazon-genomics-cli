package context

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestManager_Destroy(t *testing.T) {
	testCases := map[string]struct {
		setupMocks  func(*testing.T) mockClients
		expectedErr error
	}{
		"destroy success": {
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				defer close(mockClients.progressStream1)
				defer close(mockClients.progressStream2)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.cdkMock.EXPECT().DestroyApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Any()).Return(mockClients.progressStream1, nil)
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
		"destroy error": {
			expectedErr: fmt.Errorf("some context error"),
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
				mockClients.cdkMock.EXPECT().DestroyApp(filepath.Join(testHomeDir, ".agc/cdk/apps/context"), gomock.Any()).Return(nil, fmt.Errorf("some context error"))
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

			err := manager.Destroy(testContextName1, false)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
