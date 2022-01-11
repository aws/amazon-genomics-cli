package cfn

import "github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
func (c Client) GetStackOutputs(stackName string) (map[string]string, error) {
	info, err := c.GetStackInfo(stackName)
	return info.Outputs, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
