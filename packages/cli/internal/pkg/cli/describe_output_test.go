package cli

import (
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/stretchr/testify/assert"
)

func TestDescribeOutput(t *testing.T) {
	tests := map[string]struct {
		output              interface{}
		expectedDescription string
	}{
		"WorkflowInstance": {
			output:              types.WorkflowInstance{},
			expectedDescription: "Output of the command has following format:\nWORKFLOWINSTANCE: ContextName Id InProject State SubmittedTime WorkflowName\n",
		},
		"Workflow": {
			output:              types.Workflow{},
			expectedDescription: "Output of the command has following format:\nWORKFLOW: Name Source TypeLanguage TypeVersion\n",
		},
		"slice of WorkflowName": {
			output:              []types.WorkflowName{},
			expectedDescription: "Output of the command has following format:\nWORKFLOWNAME: Name\n",
		},
		"ContextInstance": {
			output:              types.ContextInstance{},
			expectedDescription: "Output of the command has following format:\nCONTEXTINSTANCE: ErrorStatus Id Info Name RunStatus RunTime StartTime\n",
		},
		"Context": {
			output: types.Context{},
			expectedDescription: "Output of the command has following format:\nCONTEXT: MaxVCpus Name RequestSpotInstances Status" +
				" StatusReason\nINSTANCETYPE: Value\nOUTPUTLOCATION: Url\n",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualDescription := DescribeOutput(tt.output)
			assert.Equal(t, tt.expectedDescription, actualDescription)
		})
	}
}
