package cdk

import (
	"os"
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

	return executeCdkCommandAndCleanupDirectory(appDir, cmdArgs, tmpDir)
}
