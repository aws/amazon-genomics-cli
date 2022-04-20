package osutils

import (
	"errors"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"io/fs"
	"os"
	"path/filepath"
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

func TestGetWalkDirFn(t *testing.T) {
	backupOsStat, backupOsIsNotExist, backupOsMkdirAll, backupOsCreate, backupIoCopy, backupOsOpen := osStat, osIsNotExist, osMkdirAll, osCreate, ioCopy, osOpen

	defer func() {
		osStat, osIsNotExist, osMkdirAll, osCreate, ioCopy, osOpen = backupOsStat, backupOsIsNotExist, backupOsMkdirAll, backupOsCreate, backupIoCopy, backupOsOpen
	}()

	var tests = map[string]struct {
		setupMocks        func(mockUtils MockUtils)
		currentPath       string
		absSourceDir      string
		absDestinationDir string
		errMessage        string
		expectedOutput    string
		expectedErr       error
		dirEntry          fs.DirEntry
	}{
		"empty directory returns nil": {
			setupMocks:        func(mocksUtils MockUtils) {},
			currentPath:       "/some/path/to/folder/file.ext",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/other/path",
			expectedOutput:    "/some/other/path/file.ext",
			expectedErr:       errors.New("file doesn't exist"),
			dirEntry:          nil,
		},
		"it skips file: '.nextflow'": {
			setupMocks:        func(mocksUtils MockUtils) {},
			currentPath:       "/some/path/to/folder/file.ext",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/path/to",
			expectedOutput:    "/some/path/to/folder/file.ext",
			expectedErr:       filepath.SkipDir,
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      ".nextflow",
			},
		},
		"it skips file: '.snakemake'": {
			setupMocks:        func(mocksUtils MockUtils) {},
			currentPath:       "/some/path/to/folder/file.ext",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/path/to",
			expectedOutput:    "/some/path/to/folder/file.ext",
			expectedErr:       filepath.SkipDir,
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      ".snakemake",
			},
		},
		"it skips file: 'work'": {
			setupMocks:        func(mocksUtils MockUtils) {},
			currentPath:       "/some/path/to/folder/file.ext",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/path/to",
			expectedOutput:    "/some/path/to/folder/file.ext",
			expectedErr:       filepath.SkipDir,
			dirEntry: TestDirEntry{
				isDir:     true,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      "work",
			},
		},
		"it skips suffix file: '.nextflow.log'": {
			setupMocks:        func(mocksUtils MockUtils) {},
			currentPath:       "/some/path/to/folder/file.ext",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/path/to",
			expectedOutput:    "/some/path/to/folder/file.ext",
			expectedErr:       filepath.SkipDir,
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      ".nextflow.log",
			},
		},
		"if it tries to open an invalid file, it fails - OsOpen": {
			setupMocks: func(mocksUtils MockUtils) {
				osOpen = mocksUtils.mockOs.Open
				var srcFile = &os.File{}
				mocksUtils.mockOs.EXPECT().Open("/some/path/to/folder/file.ext").Return(srcFile, errors.New("error here"))
			},
			currentPath:       "/some/path/to/folder/file.ext",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/other/path",
			expectedOutput:    "/some/path/to/folder/file.ext",
			expectedErr:       errors.New("error here"),
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      "some dir here",
			},
		},
		"if its a real directory - returns nil": {
			setupMocks:        func(mocksUtils MockUtils) {},
			currentPath:       "/some/path/to/folder",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/other/path",
			expectedOutput:    "/some/other/path",
			expectedErr:       nil,
			dirEntry: TestDirEntry{
				isDir:     true,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      "some dir here",
			},
		},
		"get relative path fails": {
			setupMocks: func(mocksUtils MockUtils) {
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osOpen = mocksUtils.mockOs.Open
				var srcFile = &os.File{}
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path").Return(nil, doesNotExist)
				mocksUtils.mockOs.EXPECT().Open("/some/path/to/folder/file.ext").Return(srcFile, nil)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true)
				returnedError := errors.New("some error")
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path", gomock.Any()).Return(returnedError)
			},
			currentPath:       "/some/path/to/folder/file.ext",
			absSourceDir:      "/some/path/to/folder",
			absDestinationDir: "/some/other/path",
			expectedOutput:    "/some/path/to/folder/file.ext",
			expectedErr:       errors.New("some other error"),
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      "",
			},
		},
		"if it tries to create a file, it fails - OsCreate": {
			setupMocks: func(mocksUtils MockUtils) {
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osCreate = mocksUtils.mockOs.Create
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				osOpen = mocksUtils.mockOs.Open
				var srcFile = &os.File{}
				var dstFile = &os.File{}
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path", gomock.Any()).Return(nil)
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path").Return(nil, doesNotExist)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true).Times(1)
				mocksUtils.mockOs.EXPECT().Open("/some/path/to/file.ext").Return(srcFile, nil)
				mocksUtils.mockOs.EXPECT().Create("/some/other/path/file.ext").Return(dstFile, errors.New("new error"))
			},
			currentPath:       "/some/path/to/file.ext",
			absSourceDir:      "/some/path/to",
			absDestinationDir: "/some/other/path",
			expectedOutput:    "/some/other/path/to/some/subfolder/file.ext",
			expectedErr:       errors.New("new error"),
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      "",
			},
		},
		"tries to create a file, it fails - IoCopy": { //err?
			setupMocks: func(mocksUtils MockUtils) {
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osCreate = mocksUtils.mockOs.Create
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				osOpen = mocksUtils.mockOs.Open
				ioCopy = mocksUtils.mockIo.Copy
				var srcFile = &os.File{}
				var dstFile = &os.File{}
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path", gomock.Any()).Return(nil).Times(1)
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path").Return(nil, doesNotExist).Times(1)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true).Times(1)
				mocksUtils.mockOs.EXPECT().Open("/some/path/to/file.ext").Return(srcFile, nil)
				mocksUtils.mockOs.EXPECT().Create("/some/other/path/file.ext").Return(dstFile, nil)
				mocksUtils.mockIo.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1), errors.New("new err"))
			},
			currentPath:       "/some/path/to/file.ext",
			absSourceDir:      "/some/path/to",
			absDestinationDir: "/some/other/path",
			expectedOutput:    "/some/other/path/to/some/subfolder/file.ext",
			expectedErr:       errors.New("new err"),
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      "",
			},
		},
		"GetWalkDirFn - happy path": {
			setupMocks: func(mocksUtils MockUtils) {
				osMkdirAll = mocksUtils.mockOs.MkdirAll
				osCreate = mocksUtils.mockOs.Create
				osStat = mocksUtils.mockOs.Stat
				osIsNotExist = mocksUtils.mockOs.IsNotExist
				osOpen = mocksUtils.mockOs.Open
				ioCopy = mocksUtils.mockIo.Copy
				var srcFile = &os.File{}
				var dstFile = &os.File{}
				doesNotExist := errors.New("Some error")
				mocksUtils.mockOs.EXPECT().MkdirAll("/some/other/path", gomock.Any()).Return(nil).Times(1)
				mocksUtils.mockOs.EXPECT().Stat("/some/other/path").Return(nil, doesNotExist).Times(1)
				mocksUtils.mockOs.EXPECT().IsNotExist(doesNotExist).Return(true).Times(1)
				mocksUtils.mockOs.EXPECT().Open("/some/path/to/file.ext").Return(srcFile, nil)
				mocksUtils.mockOs.EXPECT().Create("/some/other/path/file.ext").Return(dstFile, nil)
				mocksUtils.mockIo.EXPECT().Copy(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			currentPath:       "/some/path/to/file.ext",
			absSourceDir:      "/some/path/to",
			absDestinationDir: "/some/other/path",
			expectedOutput:    "/some/other/path/to/some/subfolder/file.ext",
			expectedErr:       nil,
			dirEntry: TestDirEntry{
				isDir:     false,
				fileMode:  32,
				info:      nil,
				infoError: nil,
				name:      "",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOs := iomocks.NewMockOS(ctrl)
			mockIo := iomocks.NewMockIO(ctrl)
			mockUtils := &MockUtils{
				ctrl:   ctrl,
				mockOs: mockOs,
				mockIo: mockIo,
			}
			tt.setupMocks(*mockUtils)

			walkDirFunc := getWalkDirFn(tt.absDestinationDir, tt.absSourceDir)
			err := walkDirFunc(tt.currentPath, tt.dirEntry, nil)

			if tt.expectedErr != nil {
				assert.Error(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type MockUtils struct {
	ctrl   *gomock.Controller
	mockOs *iomocks.MockOS
	mockIo *iomocks.MockIO
}

type TestDirEntry struct {
	isDir     bool
	fileMode  fs.FileMode
	info      fs.FileInfo
	infoError error
	name      string
}

func (tde TestDirEntry) IsDir() bool {
	return tde.isDir
}

func (tde TestDirEntry) Type() fs.FileMode {
	return tde.fileMode
}

func (tde TestDirEntry) Info() (fs.FileInfo, error) {
	return tde.info, tde.infoError
}

func (tde TestDirEntry) Name() string {
	return tde.name
}
