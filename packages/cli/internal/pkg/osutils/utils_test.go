package osutils

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDetermineHomeDir_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osUserHomeDir = mockOs.UserHomeDir
	expectedPath := "/some/dir"
	mockOs.EXPECT().UserHomeDir().Return(expectedPath, nil)
	actualPath, err := DetermineHomeDir()

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)
}

func TestDetermineHomeDir_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osUserHomeDir = mockOs.UserHomeDir
	expectedOsError := errors.New("some error")
	mockOs.EXPECT().UserHomeDir().Return("", expectedOsError)
	_, err := DetermineHomeDir()

	expectedError := actionableerror.New(err, "Please check that your home or user profile directory is defined within your environment variables")

	assert.Error(t, err, expectedError)
}

func TestExpandHomeDir_WithExpansion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osUserHomeDir = mockOs.UserHomeDir

	expectedHomePath := "home/user"
	expectedExpandedPath := "home/user/FooBar"

	mockOs.EXPECT().UserHomeDir().Return(expectedHomePath, nil)
	actualExpandedPath := ExpandHomeDir("~/FooBar")

	assert.Equal(t, expectedExpandedPath, actualExpandedPath)
	ctrl.Finish()
}

func TestExpandHomeDir_WithoutExpansion(t *testing.T) {
	expectedExpandedPath := "FooBar"
	actualExpandedPath := ExpandHomeDir("FooBar")

	assert.Equal(t, expectedExpandedPath, actualExpandedPath)
}
