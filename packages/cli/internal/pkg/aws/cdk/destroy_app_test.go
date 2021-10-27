package cdk

import (
	"testing"

	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testDestroyAppPath         = "test/destroy/app/path"
	testDestroyProfile         = "test-destroy-profile"
	testDestroyEnvironmentVars = []string{"a=1", "b=2"}
	testDestroyUniqueKey       = "test-unique-key"
)

func TestDestroyApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	mkDirTemp = mockOs.MkdirTemp
	client := NewClient(testDestroyProfile)
	realExecuteCdkCommand := ExecuteCdkCommand
	ExecuteCdkCommand = func(cdkAppPathDir string, cmdArgs []string, uniqueKey string) (ProgressStream, error) {
		assert.Equal(t, testDestroyAppPath, cdkAppPathDir)
		assert.Equal(t, []string{
			"destroy",
			"--all",
			"--force",
			"--profile", testDestroyProfile,
			"--output", "/some/path",
			"-c", "a=1",
			"-c", "b=2",
		}, cmdArgs)
		return nil, nil
	}
	defer func() { ExecuteCdkCommand = realExecuteCdkCommand }()
	mockOs.EXPECT().MkdirTemp(testDestroyAppPath, "cdk-output").Return("/some/path", nil)

	_, _ = client.DestroyApp(testDestroyAppPath, testDestroyEnvironmentVars, testDestroyUniqueKey)
}

func TestDestroyAppWithNilEnvVars(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	mkDirTemp = mockOs.MkdirTemp
	client := NewClient("")
	realExecuteCdkCommand := ExecuteCdkCommand
	ExecuteCdkCommand = func(cdkAppPathDir string, cmdArgs []string, uniqueKey string) (ProgressStream, error) {
		assert.Equal(t, testDestroyAppPath, cdkAppPathDir)
		assert.Equal(t, []string{
			"destroy",
			"--all",
			"--force",
			"--profile", "",
			"--output", "/some/path",
		}, cmdArgs)
		return nil, nil
	}
	defer func() { ExecuteCdkCommand = realExecuteCdkCommand }()
	mockOs.EXPECT().MkdirTemp(testDestroyAppPath, "cdk-output").Return("/some/path", nil)

	_, _ = client.DestroyApp(testDestroyAppPath, nil, testDestroyUniqueKey)
}
