package cli

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	workflowmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/workflow"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	testWorkflowName1 = "TestWorkflow1"
	testWorkflowName2 = "TestWorkflow2"
)

func TestWorkflowAutoComplete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowMock := workflowmocks.NewMockWorkflowManager(ctrl)
	workflowMock.EXPECT().ListWorkflows().Return(map[string]workflow.Summary{testWorkflowName1: {}, testWorkflowName2: {}}, nil)
	workflowAutoComplete := &WorkflowAutoComplete{
		workflowManagerFactory: func() workflow.Interface {
			return workflowMock
		},
	}
	autoCompleteFunction := workflowAutoComplete.GetWorkflowAutoComplete()
	_, compDirective := autoCompleteFunction(nil, make([]string, 0), "")
	assert.Equal(t, compDirective, cobra.ShellCompDirectiveNoFileComp)
}

func TestWorkflowAutoComplete_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowMock := workflowmocks.NewMockWorkflowManager(ctrl)
	workflowMock.EXPECT().ListWorkflows().Return(map[string]workflow.Summary{}, errors.New("Test Workflow Error"))
	workflowAutoComplete := &WorkflowAutoComplete{
		workflowManagerFactory: func() workflow.Interface {
			return workflowMock
		},
	}
	autoCompleteFunction := workflowAutoComplete.GetWorkflowAutoComplete()
	actualKeys, compDirective := autoCompleteFunction(nil, make([]string, 0), "")
	assert.Equal(t, []string(nil), actualKeys)
	assert.Equal(t, compDirective, cobra.ShellCompDirectiveError)
}
