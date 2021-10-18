package context

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/stretchr/testify/assert"
)

func TestManager_List(t *testing.T) {
	testCases := map[string]struct {
		setupMocks       func(*testing.T) mockClients
		expectedContexts map[string]Summary
		expectedErr      error
	}{
		"contexts": {
			expectedContexts: map[string]Summary{
				testContextName1: {
					Name: testContextName1,
					Engines: []spec.Engine{
						{
							Type:   "wdl",
							Engine: "cromwell",
						},
					},
				},
				testContextName2: {
					Name: testContextName2,
					Engines: []spec.Engine{
						{
							Type:   "wdl",
							Engine: "cromwell",
						},
					},
				},
				testContextName3: {
					Name: testContextName3,
					Engines: []spec.Engine{
						{
							Type:   "wdl",
							Engine: "cromwell",
						},
					},
				},
			},
			setupMocks: func(t *testing.T) mockClients {
				mockClients := createMocks(t)
				mockClients.configMock.EXPECT().GetUserEmailAddress().Return(testUserEmail, nil)
				mockClients.configMock.EXPECT().GetUserId().Return(testUserId, nil)
				mockClients.projMock.EXPECT().Read().Return(testValidProjectSpec, nil)
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

			contexts, err := manager.List()

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedContexts, contexts)
			}

		})
	}
}
