package cdk

import "github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"

func (client Client) DestroyApp(appDir string, context []string) (ProgressStream, error) {
	cmdArgs := []string{
		"destroy",
		"--all",
		"--force",
		"--profile", client.profile,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)

	progressStream, err := ExecuteCdkCommand(appDir, cmdArgs)
	return progressStream, actionable.FindSuggestionForError(err, actionable.AwsErrorMessageToSuggestedActionMap)
}
