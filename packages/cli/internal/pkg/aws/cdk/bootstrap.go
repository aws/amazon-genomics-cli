package cdk

import (
	"fmt"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/awsresources"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
)

func (client Client) Bootstrap(appDir string, context []string, executionName string) (ProgressStream, error) {
	cmdArgs := []string{
		"bootstrap",
		"--toolkit-stack-name", awsresources.RenderBootstrapStackName(),
		"--qualifier", constants.AppTagValue,
		"--tags", fmt.Sprintf("%s=%s", constants.AppTagKey, constants.AppTagValue),
		"--profile", client.profile,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)
	progressStream, err := executeCdkCommand(appDir, cmdArgs, executionName)
	return progressStream, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
