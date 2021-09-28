package cfn

import "github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"

func (c Client) GetStackOutputs(stackName string) (map[string]string, error) {
	info, err := c.GetStackInfo(stackName)
	return info.Outputs, actionable.FindSuggestionForError(err, actionable.AwsErrorMessageToSuggestedActionMap)
}
