package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	managermocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_workflowOutputOpts_Validate(t *testing.T) {

	tests := []struct {
		name    string
		vars    workflowOutputVars
		wantErr bool
	}{
		{
			name:    "validate valid input",
			vars:    workflowOutputVars{runId: "abcd"},
			wantErr: false,
		},
		{
			name:    "invalidate invalid input",
			vars:    workflowOutputVars{runId: "   "},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &workflowOutputOpts{
				vars: tt.vars,
			}
			if err := o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkflowOutputOpts_Execute(t *testing.T) {

	const (
		testInstanceId1 = "Test Instance Id 1"
		testInstanceId2 = "Test Instance Id 2"
	)
	testInstanceId2Err := actionableerror.New(fmt.Errorf("no workflow run found for workflow run '%s'", testInstanceId2), "check the workflow run id and check the workflow was run from the current project")

	tests := map[string]struct {
		setupOpts      func(opts *workflowOutputOpts)
		expectedOutput map[string]interface{}
		expectedErr    error
	}{
		"instanceIdFound": {
			setupOpts: func(opts *workflowOutputOpts) {
				opts.vars.runId = testInstanceId1
				opts.wfManager.(*managermocks.MockWorkflowManager).EXPECT().OutputByInstanceId(testInstanceId1).
					Times(1).
					Return(map[string]interface{}{"foo": "baa"}, nil)
			},
			expectedOutput: map[string]interface{}{"foo": "baa"},
			expectedErr:    nil,
		},
		"instanceIdNotFound": {
			setupOpts: func(opts *workflowOutputOpts) {
				opts.vars.runId = testInstanceId2
				opts.wfManager.(*managermocks.MockWorkflowManager).EXPECT().OutputByInstanceId(testInstanceId2).
					Times(1).
					Return(nil, actionableerror.New(fmt.Errorf("no workflow run found for workflow run '%s'", testInstanceId2), "check the workflow run id and check the workflow was run from the current project"))
			},
			expectedOutput: nil,
			expectedErr:    testInstanceId2Err,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockManager := managermocks.NewMockWorkflowManager(ctrl)
			opts := &workflowOutputOpts{
				wfManager: mockManager,
			}
			tt.setupOpts(opts)
			actualOutput, err := opts.Execute()

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, actualOutput)
			} else {
				assert.EqualErrorf(t, err, tt.expectedErr.Error(), "expected error message '%v', but got '%v'", tt.expectedErr.Error(), err)
			}
		})
	}
}
