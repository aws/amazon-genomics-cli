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
		input           string
		mockCallTimes   int
		expectedHomeDir string
		err             error
		expectedOutput  string
	}{
		"HomeDir only": {
			input:           "~",
			expectedHomeDir: "/home/user",
			err:             nil,
			mockCallTimes:   1,
			expectedOutput:  "/home/user",
		},
		"HomeDir only with backslash": {
			input:           "~/",
			expectedHomeDir: "/home/user",
			err:             nil,
			mockCallTimes:   1,
			expectedOutput:  "/home/user",
		},
		"HomeDir only with subsequent dir": {
			input:           "~/FooBar",
			expectedHomeDir: "/home/user",
			err:             nil,
			mockCallTimes:   1,
			expectedOutput:  "/home/user/FooBar",
		},
		"Tilde at the beginning": {
			input:           "~~~",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "~~~",
		},
		"Tilde at the beginning with backslash": {
			input:           "/~/~/~/",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "/~/~/~/",
		},
		"Tilde at the beginning ": {
			input:           "/~/~/~",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "/~/~/~",
		},
		"Tilde in the middle": {
			input:           "Foo~/~Bar",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "Foo~/~Bar",
		},
		"Tilde in the end": {
			input:           "~Foo/Bar~",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "~Foo/Bar~",
		},
		"Empty string": {
			input:           "",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "",
		},
		"Relative": {
			input:           "FooBar",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "FooBar",
		},
		"Absolute": {
			input:           "/Foo/Bar",
			expectedHomeDir: "",
			err:             nil,
			mockCallTimes:   0,
			expectedOutput:  "/Foo/Bar",
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

func TestGetAndCreateRelativePath(t *testing.T) {
	backupOsStat, backupOsIsNotExist, backupOsMkdirAll := osStat, osIsNotExist, osMkdirAll

	defer func() {
		osStat, osIsNotExist, osMkdirAll = backupOsStat, backupOsIsNotExist, backupOsMkdirAll
	}()

	var tests = map[string]struct {
		setupMocks     func(mockUtils MockUtils)
		currentPath    string
		sourcePath     string
		destinationDir string
		errMessage     string
		expectedOutput string
	}{
		"file at root": {
			setupMocks: func(mocksUtils MockUtils) {
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path", gomock.Any()).Return(nil).Times(1)
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path").Return(nil, doesNotExist).Times(1)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true).Times(1)
			},
			currentPath:    "/some/path/to/folder/file.ext",
			sourcePath:     "/some/path/to/folder",
			destinationDir: "/some/other/path",
			errMessage:     "",
			expectedOutput: "/some/other/path/file.ext",
		},
		"file one level deep": {
			setupMocks: func(mocksUtils MockUtils) {
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path/subfolder").Return(nil, doesNotExist).Times(1)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true).Times(1)
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path/subfolder", gomock.Any()).Return(nil).Times(1)
			},
			currentPath:    "/some/path/to/folder/subfolder/file.ext",
			sourcePath:     "/some/path/to/folder",
			destinationDir: "/some/other/path",
			errMessage:     "",
			expectedOutput: "/some/other/path/subfolder/file.ext",
		},
		"Duplicate sub paths aren't replaced": {
			setupMocks: func(mocksUtils MockUtils) {
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path/some/path/to/folder/subfolder").Return(nil, doesNotExist).Times(1)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true).Times(1)
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path/some/path/to/folder/subfolder", gomock.Any()).Return(nil).Times(1)
			},
			currentPath:    "/some/path/to/folder/some/path/to/folder/subfolder/file.ext",
			sourcePath:     "/some/path/to/folder",
			destinationDir: "/some/other/path",
			errMessage:     "",
			expectedOutput: "/some/other/path/some/path/to/folder/subfolder/file.ext",
		},
		"fails to create folder": {
			setupMocks: func(mocksUtils MockUtils) {
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path/subfolder").Return(nil, doesNotExist).Times(1)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true).Times(1)
				returnedError := errors.New("failed to create")
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path/subfolder", gomock.Any()).Return(returnedError).Times(1)
			},
			currentPath:    "/some/path/to/folder/subfolder/file.ext",
			sourcePath:     "/some/path/to/folder",
			destinationDir: "/some/other/path",
			errMessage:     "failed to create",
			expectedOutput: "/some/other/path/subfolder/file.ext",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOs := iomocks.NewMockOS(ctrl)
			mockUtils := &MockUtils{
				ctrl:   ctrl,
				mockOs: mockOs,
			}
			tt.setupMocks(*mockUtils)

			actualPath, err := getAndCreateRelativePath(tt.currentPath, tt.sourcePath, tt.destinationDir)
			if err != nil {
				assert.Error(t, err, tt.errMessage)
			} else {
				assert.Equal(t, tt.expectedOutput, actualPath)
			}
		})
	}
}

type MockUtils struct {
	ctrl   *gomock.Controller
	mockOs *iomocks.MockOS
}
