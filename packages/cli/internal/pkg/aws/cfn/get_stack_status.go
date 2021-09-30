package cfn

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func (c Client) GetStackStatus(stackName string) (types.StackStatus, error) {
	info, err := c.GetStackInfo(stackName)
	return info.Status, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
