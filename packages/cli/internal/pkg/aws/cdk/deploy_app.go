package cdk

import (
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/awsresources"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
)

var mkDirTemp = os.MkdirTemp

func (client Client) DeployApp(appDir string, context []string, executionName string) (ProgressStream, error) {
	tmpDir, _ := mkDirTemp(appDir, "cdk-output")
	// The temp directory defaults to 700 permissions, which can prevent Docker
	// from mounting directories under it, which can prevent the CDK's bundling
	// process from working, with errors like:
	// docker: Error response from daemon: error while creating mount source path '...': mkdir ...: permission denied.
	// So we need to change the permissions
	err := os.Chmod(tmpDir, 0755)
	if err != nil {
		return nil, err
	}
	cmdArgs := []string{
		"deploy",
		"--all",
		"--profile", client.profile,
		"--require-approval", "never",
		"--toolkit-stack-name", awsresources.RenderBootstrapStackName(),
		"--output", tmpDir,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)
	progressStream, err := executeCdkCommandAndCleanupDirectory(appDir, cmdArgs, tmpDir, executionName)
	return progressStream, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
}
