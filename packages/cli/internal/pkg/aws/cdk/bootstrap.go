package cdk

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
)

func (client Client) Bootstrap(appDir string, context []string, executionName string) (ProgressStream, error) {
	cmdArgs := []string{
		"bootstrap",
		"--profile", client.profile,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)
	progressStream, err := executeCdkCommand(appDir, cmdArgs, executionName)
	return progressStream, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
