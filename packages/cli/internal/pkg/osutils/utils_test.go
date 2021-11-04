package osutils

import (
	"errors"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
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
	mockUtils := iomocks.NewMockUtils(ctrl)

	DetermineHomeDir = func() { return "home/user"}

	expectedHomePath := "home/user"
	expectedExpandedPath := "home/user/FooBar"

	mockExpandHomeDir.EXPECT().DetermineHomeDir().Return(expectedHomePath, nil).Times(1)
	actualExpandedPath := ExpandHomeDir("~/FooBar")

	assert.Equal(t, expectedExpandedPath, actualExpandedPath)
	ctrl.Finish()
}
//
//func TestExpandHomeDir_WithoutExpansion(t *testing.T) {
//	//ctrl := gomock.NewController(t)
//	//mockExpandHomeDir := iomocks.NewMockExpandHomeDir(ctrl)
//
//	expectedExpandedPath := "some/dir/FooBar"
//
//	//mockExpandHomeDir.EXPECT().DetermineHomeDir().Return(expectedHomePath, nil).Times(1)
//	actualExpandedPath := ExpandHomeDir("some/dir/FooBar")
//
//	assert.Equal(t, expectedExpandedPath, actualExpandedPath)
//	//ctrl.Finish()
//}
