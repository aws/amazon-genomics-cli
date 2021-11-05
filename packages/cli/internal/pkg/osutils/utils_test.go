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

	expectedHomePath := "/home/user"
	expectedExpandedPath := "/home/user/FooBar"

	mockOs.EXPECT().UserHomeDir().Return(expectedHomePath, nil)
	actualExpandedPath, _ := ExpandHomeDir("~/FooBar")

	assert.Equal(t, expectedExpandedPath, actualExpandedPath)
	ctrl.Finish()
}

func TestExpandHomeDir_WithExpansion_HomeDirOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osUserHomeDir = mockOs.UserHomeDir

	expectedHomePath := "/home/user"
	expectedExpandedPath := "/home/user"

	mockOs.EXPECT().UserHomeDir().Return(expectedHomePath, nil)
	actualExpandedPath, _ := ExpandHomeDir("~")

	assert.Equal(t, expectedExpandedPath, actualExpandedPath)
	ctrl.Finish()
}

func TestExpandHomeDir_WithExpansion_HomeDirOnlyWithSlash(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osUserHomeDir = mockOs.UserHomeDir

	expectedHomePath := "/home/user"
	expectedExpandedPath := "/home/user"

	mockOs.EXPECT().UserHomeDir().Return(expectedHomePath, nil)
	actualExpandedPath, _ := ExpandHomeDir("~/")

	assert.Equal(t, expectedExpandedPath, actualExpandedPath)
	ctrl.Finish()
}

func TestExpandHomeDir_WithExpansionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osUserHomeDir = mockOs.UserHomeDir
	expectedOsError := errors.New("OS Error")
	mockOs.EXPECT().UserHomeDir().Return("", expectedOsError)
	_, err := ExpandHomeDir("~/Error/Dir")

	expectedError := actionableerror.New(err, "Please check that your home or user profile directory is defined within your environment variables")

	assert.Error(t, err, expectedError)
	ctrl.Finish()
}

func TestExpandHomeDir_WithoutExpansion(t *testing.T) {
	tests := map[string]struct {
		input        string
		expectedPath string
	}{
		"Tilde at the beginning": {
			input:        "~~~",
			expectedPath: "~~~",
		},
		"Tilde at the beginning with backslash": {
			input:        "/~/~/~/",
			expectedPath: "/~/~/~/",
		},
		"Tilde at the beginning ": {
			input:        "/~/~/~",
			expectedPath: "/~/~/~",
		},
		"Tilde in the middle": {
			input:        "Foo~/~Bar",
			expectedPath: "Foo~/~Bar",
		},
		"Tilde in the end": {
			input:        "~Foo/Bar~",
			expectedPath: "~Foo/Bar~",
		},
		"Empty string": {
			input:        "",
			expectedPath: "",
		},
		"Relative": {
			input:        "FooBar",
			expectedPath: "FooBar",
		},
		"Absolute": {
			input:        "/Foo/Bar",
			expectedPath: "/Foo/Bar",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualPath, _ := ExpandHomeDir(tt.input)
			assert.Equal(t, tt.expectedPath, actualPath)
		})
	}
}
