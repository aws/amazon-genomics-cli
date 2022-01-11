package cli

import (
	"fmt"
	"testing"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	managermocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStatusWorkflowOpts_Validate(t *testing.T) {
	tests := map[int]bool{
		-1:     false,
		0:      false,
		1:      true,
		2:      true,
		10:     true,
		100:    true,
		999:    true,
		1000:   true,
		1001:   false,
		100000: false,
	}

	for maxInstances, isValid := range tests {
		t.Run(fmt.Sprint(maxInstances), func(t *testing.T) {
			opts := &workflowStatusOpts{
				workflowStatusVars: workflowStatusVars{MaxInstances: maxInstances},
			}
			err := opts.Validate()
			if isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestWorkflowStatusOpts_Execute(t *testing.T) {

	const (
		testContext      = "Test Context"
		testMaxInstances = 34

		testInstanceId1 = "Test Instance Id 1"
		testWorkflow1   = "Test Workflow 1"
		testSubmitTime1 = "Test Submit Time 1"
		testState1      = "Test State 1"

		testInstanceId2 = "Test Instance Id 2"
		testWorkflow2   = "Test Workflow 2"
		testSubmitTime2 = "Test Submit Time 2"
		testState2      = "Test State 2"
	)

	testInstanceSummary1 := workflow.InstanceSummary{
		Id:           testInstanceId1,
		WorkflowName: testWorkflow1,
		ContextName:  testContext,
		SubmitTime:   testSubmitTime1,
		State:        testState1,
		InProject:    true,
	}

	testWorkflowInstance1 := types.WorkflowInstance{
		Id:            testInstanceId1,
		WorkflowName:  testWorkflow1,
		ContextName:   testContext,
		State:         testState1,
		InProject:     true,
		SubmittedTime: testSubmitTime1,
	}

	testInstanceSummary2 := workflow.InstanceSummary{
		Id:           testInstanceId2,
		WorkflowName: testWorkflow2,
		ContextName:  testContext,
		SubmitTime:   testSubmitTime2,
		State:        testState2,
		InProject:    true,
	}

	testWorkflowInstance2 := types.WorkflowInstance{
		Id:            testInstanceId2,
		WorkflowName:  testWorkflow2,
		ContextName:   testContext,
		State:         testState2,
		InProject:     true,
		SubmittedTime: testSubmitTime2,
	}

	tests := map[string]struct {
		setupOpts        func(*workflowStatusOpts)
		expectedStatuses []types.WorkflowInstance
	}{
		"byInstanceId": {
			setupOpts: func(opts *workflowStatusOpts) {
				opts.InstanceId = testInstanceId1
				opts.wfManager.(*managermocks.MockWorkflowManager).EXPECT().StatusWorkflowByInstanceId(testInstanceId1).
					Times(1).
					Return([]workflow.InstanceSummary{testInstanceSummary1}, nil)
			},
			expectedStatuses: []types.WorkflowInstance{testWorkflowInstance1},
		},
		"byWorkflowName": {
			setupOpts: func(opts *workflowStatusOpts) {
				opts.WorkflowName = testWorkflow1
				opts.MaxInstances = testMaxInstances
				opts.wfManager.(*managermocks.MockWorkflowManager).EXPECT().StatusWorkflowByName(testWorkflow1, testMaxInstances).
					Times(1).
					Return([]workflow.InstanceSummary{testInstanceSummary1}, nil)
			},
			expectedStatuses: []types.WorkflowInstance{testWorkflowInstance1},
		},
		"byContextName": {
			setupOpts: func(opts *workflowStatusOpts) {
				opts.ContextName = testContext
				opts.MaxInstances = testMaxInstances
				opts.wfManager.(*managermocks.MockWorkflowManager).EXPECT().StatusWorkflowByContext(testContext, testMaxInstances).
					Times(1).
					Return([]workflow.InstanceSummary{testInstanceSummary1, testInstanceSummary2}, nil)
			},
			expectedStatuses: []types.WorkflowInstance{testWorkflowInstance1, testWorkflowInstance2},
		},
		"Default": {
			setupOpts: func(opts *workflowStatusOpts) {
				opts.MaxInstances = testMaxInstances
				opts.wfManager.(*managermocks.MockWorkflowManager).EXPECT().StatusWorkflowAll(testMaxInstances).
					Times(1).
					Return([]workflow.InstanceSummary{testInstanceSummary1, testInstanceSummary2}, nil)
			},
			expectedStatuses: []types.WorkflowInstance{testWorkflowInstance1, testWorkflowInstance2},
		},
		"Default_EmptyResult": {
			setupOpts: func(opts *workflowStatusOpts) {
				opts.MaxInstances = testMaxInstances
				opts.wfManager.(*managermocks.MockWorkflowManager).EXPECT().StatusWorkflowAll(testMaxInstances).
					Times(1).
					Return([]workflow.InstanceSummary{}, nil)
			},
			expectedStatuses: []types.WorkflowInstance{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockManager := managermocks.NewMockWorkflowManager(ctrl)
			opts := &workflowStatusOpts{
				wfManager: mockManager,
			}
			tt.setupOpts(opts)
			actualStatuses, err := opts.Execute()

			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectedStatuses, actualStatuses)
			}
		})
	}
}
