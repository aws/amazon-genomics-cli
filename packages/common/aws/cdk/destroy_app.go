package cdk

func (client Client) DestroyApp(appDir string, context []string) (ProgressStream, error) {
	cmdArgs := []string{
		"destroy",
		"--all",
		"--force",
		"--profile", client.profile,
	}
	cmdArgs = appendContextArguments(cmdArgs, context)

	return ExecuteCdkCommand(appDir, cmdArgs)
}
