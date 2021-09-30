package cdk

import "github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"

func (client Client) DestroyApp(appDir string, context []string) (ProgressStream, error) {
	cmdArgs := []string{
		"destroy",
		"--all",
		"--force",
		"--profile", client.profile,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)

	progressStream, err := ExecuteCdkCommand(appDir, cmdArgs)
	return progressStream, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
