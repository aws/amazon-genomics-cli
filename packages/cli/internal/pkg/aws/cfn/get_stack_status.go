package cfn

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func (c Client) GetStackStatus(stackName string) (types.StackStatus, error) {
	info, err := c.GetStackInfo(stackName)
	return info.Status, actionable.FindSuggestionForError(err, actionable.AwsErrorMessageToSuggestedActionMap)
}
