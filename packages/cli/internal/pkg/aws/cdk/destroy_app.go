package cdk

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
)

func (client Client) DestroyApp(appDir string, context []string, deploymentName string) (ProgressStream, error) {
	tmpDir, _ := mkDirTemp(appDir, "cdk-output")
	cmdArgs := []string{
		"destroy",
		"--all",
		"--force",
		"--profile", client.profile,
		"--output", tmpDir,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)
	progressStream, err := executeCdkCommandAndCleanupDirectory(appDir, cmdArgs, tmpDir, deploymentName)
	return progressStream, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
