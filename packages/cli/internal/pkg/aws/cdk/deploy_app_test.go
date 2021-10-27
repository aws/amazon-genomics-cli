package cdk

import (
	"testing"

	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testDeployAppPath         = "test/deploy/app/path"
	testDeployProfile         = "test-deploy-profile"
	testDeployUniqueKey       = "test-unique-key"
	testDeployEnvironmentVars = []string{"a=1", "b=2"}
)

func TestDeployApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	mkDirTemp = mockOs.MkdirTemp

	client := NewClient(testDeployProfile)
	realExecuteCdkCommand := ExecuteCdkCommand
	ExecuteCdkCommand = func(appDir string, cmdArgs []string, uniqueKey string) (ProgressStream, error) {
		assert.Equal(t, testDeployAppPath, appDir)
		assert.Equal(t, []string{
			"deploy",
			"--all",
			"--profile", testDeployProfile,
			"--require-approval", "never",
			"--output", "/some/path",
			"-c", "a=1",
			"-c", "b=2",
		}, cmdArgs)
		return nil, nil
	}
	defer func() { ExecuteCdkCommand = realExecuteCdkCommand }()

	mockOs.EXPECT().MkdirTemp(testDeployAppPath, "cdk-output").Return("/some/path", nil)

	_, _ = client.DeployApp(testDeployAppPath, testDeployEnvironmentVars, testDeployUniqueKey)
}

func TestDeployAppWithNilEnvVars(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	mkDirTemp = mockOs.MkdirTemp
	client := NewClient("")
	realExecuteCdkCommand := ExecuteCdkCommand
	ExecuteCdkCommand = func(appDir string, cmdArgs []string, uniqueKey string) (ProgressStream, error) {
		assert.Equal(t, testDeployAppPath, appDir)
		assert.Equal(t, []string{
			"deploy",
			"--all",
			"--profile", "",
			"--require-approval", "never",
			"--output", "/some/path",
		}, cmdArgs)
		return nil, nil
	}
	defer func() { ExecuteCdkCommand = realExecuteCdkCommand }()

	mockOs.EXPECT().MkdirTemp(testDeployAppPath, "cdk-output").Return("/some/path", nil)

	_, _ = client.DeployApp(testDeployAppPath, nil, testDeployUniqueKey)
}
