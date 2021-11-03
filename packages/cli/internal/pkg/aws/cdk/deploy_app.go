package cdk

import (
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
)

var mkDirTemp = os.MkdirTemp

func (client Client) DeployApp(appDir string, context []string) (ProgressStream, error) {
	tmpDir, _ := mkDirTemp(appDir, "cdk-output")
	cmdArgs := []string{
		"deploy",
		"--all",
		"--profile", client.profile,
		"--require-approval", "never",
		"--output", tmpDir,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)
	progressStream, err := executeCdkCommandAndCleanupDirectory(appDir, cmdArgs, tmpDir)
	return progressStream, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
