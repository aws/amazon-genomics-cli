package cdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDestroyAppPath         = "test/destroy/app/path"
	testDestroyProfile         = "test-destroy-profile"
	testDestroyEnvironmentVars = []string{"a=1", "b=2"}
)

func TestDestroyApp(t *testing.T) {
	client := NewClient(testDestroyProfile)
	realExecuteCdkCommand := ExecuteCdkCommand
	ExecuteCdkCommand = func(cdkAppPathDir string, cmdArgs []string) (ProgressStream, error) {
		assert.Equal(t, testDestroyAppPath, cdkAppPathDir)
		assert.Equal(t, []string{
			"destroy",
			"--all",
			"--force",
			"--profile", testDestroyProfile,
			"-c", "a=1",
			"-c", "b=2",
		}, cmdArgs)
		return nil, nil
	}
	defer func() { ExecuteCdkCommand = realExecuteCdkCommand }()

	_, _ = client.DestroyApp(testDestroyAppPath, testDestroyEnvironmentVars)
}

func TestDestroyAppWithNilEnvVars(t *testing.T) {
	client := NewClient("")
	realExecuteCdkCommand := ExecuteCdkCommand
	ExecuteCdkCommand = func(cdkAppPathDir string, cmdArgs []string) (ProgressStream, error) {
		assert.Equal(t, testDestroyAppPath, cdkAppPathDir)
		assert.Equal(t, []string{
			"destroy",
			"--all",
			"--force",
			"--profile", "",
		}, cmdArgs)
		return nil, nil
	}
	defer func() { ExecuteCdkCommand = realExecuteCdkCommand }()

	_, _ = client.DestroyApp(testDestroyAppPath, nil)
}
