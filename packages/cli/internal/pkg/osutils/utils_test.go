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

func TestExpandHomeDirNegative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOs := iomocks.NewMockOS(ctrl)
	osUserHomeDir = mockOs.UserHomeDir
	expectedOsError := errors.New("OS Error")
	mockOs.EXPECT().UserHomeDir().Return("", expectedOsError)
	_, err := ExpandHomeDir("~/Error/Dir")

	expectedError := actionableerror.New(err, "Please check that your home or user profile directory is defined within your environment variables")

	assert.Error(t, err, expectedError)
}

func TestExpandHomeDirPositive(t *testing.T) {
	var tests = map[string]struct {
		input string
		mockCallTimes int
		expectedHomeDir string
		err error
		expectedOutput string
	}{
		"HomeDir only": {
			input:               "~",
			expectedHomeDir:     "/home/user",
			err:                 nil,
			mockCallTimes:       1,
			expectedOutput:      "/home/user",
		},
		"HomeDir only with backslash": {
			input:              "~/",
			expectedHomeDir:    "/home/user",
			err:                 nil,
			mockCallTimes:       1,
			expectedOutput:      "/home/user",
		},
		"HomeDir only with subsequent dir": {
			input:              "~/FooBar",
			expectedHomeDir:    "/home/user",
			err:                nil,
			mockCallTimes:      1,
			expectedOutput:     "/home/user/FooBar",
		},
		"Tilde at the beginning": {
			input: 				"~~~",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "~~~",
		},
		"Tilde at the beginning with backslash": {
			input:              "/~/~/~/",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "/~/~/~/",
		},
		"Tilde at the beginning ": {
			input:              "/~/~/~",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "/~/~/~",
		},
		"Tilde in the middle": {
			input:              "Foo~/~Bar",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "Foo~/~Bar",
		},
		"Tilde in the end": {
			input:              "~Foo/Bar~",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "~Foo/Bar~",
		},
		"Empty string": {
			input:              "",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "",
		},
		"Relative": {
			input:              "FooBar",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "FooBar",
		},
		"Absolute": {
			input:              "/Foo/Bar",
			expectedHomeDir:    "",
			err:                nil,
			mockCallTimes:      0,
			expectedOutput:     "/Foo/Bar",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOs := iomocks.NewMockOS(ctrl)
			osUserHomeDir = mockOs.UserHomeDir
			mockOs.EXPECT().UserHomeDir().Return(tt.expectedHomeDir, tt.err).Times(tt.mockCallTimes)
			actualPath, _ := ExpandHomeDir(tt.input)
			assert.Equal(t, tt.expectedOutput, actualPath)
		})
	}
}