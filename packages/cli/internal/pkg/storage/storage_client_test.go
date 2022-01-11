// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFileName = "path/test-file.txt"
	testFileUrl  = "file://path/test-file.txt"
	testData     = "test-data"
)

func TestNewStorageInstance(t *testing.T) {
	testCases := map[string]struct {
		fileSystem  afero.Fs
		expectedErr error
	}{
		"storage instance with provided file system": {
			fileSystem: afero.NewMemMapFs(),
		},
		"storage instance without provided file system": {},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			if tc.fileSystem != nil {
				storageInstance, _ := NewStorageInstance(tc.fileSystem)
				assert.Equal(t, tc.fileSystem, storageInstance.fsUtils.Fs)
			} else {
				storageInstance, _ := NewStorageInstance()
				assert.IsType(t, &afero.OsFs{}, storageInstance.fsUtils.Fs)
			}
		})
	}
}

func TestStorageInstance_ReadAsBytes(t *testing.T) {
	testCases := map[string]struct {
		url                   string
		setupFileSystem       func() *afero.Afero
		expectedBytesAsString string
		expectedErr           error
	}{
		"read from filename": {
			url:                   testFileName,
			expectedBytesAsString: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				_ = afero.WriteFile(testFS, testFileName, []byte(testData), fs.ModePerm)
				return &afero.Afero{Fs: testFS}
			},
		},
		"read from file url": {
			url:                   testFileUrl,
			expectedBytesAsString: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				_ = afero.WriteFile(testFS, testFileName, []byte(testData), fs.ModePerm)
				return &afero.Afero{Fs: testFS}
			},
		},
		"read error": {
			url:         testFileName,
			expectedErr: fmt.Errorf("couldn't read file %s: open %s: file does not exist", testFileName, testFileName),
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: testFS}
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fsUtils := tc.setupFileSystem()
			si := StorageInstance{fsUtils}

			fileBytes, err := si.ReadAsBytes(tc.url)

			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expectedBytesAsString, string(fileBytes))
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
				require.Nil(t, fileBytes)
			}
		})
	}
}

func TestStorageInstance_ReadAsString(t *testing.T) {
	testCases := map[string]struct {
		url             string
		setupFileSystem func() *afero.Afero
		expectedString  string
		expectedErr     error
	}{
		"read from filename": {
			url:            testFileName,
			expectedString: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				_ = testFS.MkdirAll(filepath.Dir(testFileName), fs.ModeDir)
				_ = afero.WriteFile(testFS, testFileName, []byte(testData), fs.ModePerm)
				return &afero.Afero{Fs: testFS}
			},
		},
		"read from file url": {
			url:            testFileUrl,
			expectedString: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				_ = testFS.MkdirAll(filepath.Dir(testFileName), fs.ModeDir)
				_ = afero.WriteFile(testFS, testFileName, []byte(testData), fs.ModePerm)
				return &afero.Afero{Fs: testFS}
			},
		},
		"read error": {
			url:         testFileName,
			expectedErr: fmt.Errorf("couldn't read file %s: open %s: file does not exist", testFileName, testFileName),
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: testFS}
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fsUtils := tc.setupFileSystem()
			si := StorageInstance{fsUtils}

			fileString, err := si.ReadAsString(tc.url)

			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
			}
			require.Equal(t, tc.expectedString, fileString)
		})
	}
}

func TestStorageInstance_WriteFromBytes(t *testing.T) {
	testCases := map[string]struct {
		url             string
		data            string
		setupFileSystem func() *afero.Afero
		expectedErr     error
	}{
		"write to filename": {
			url:  testFileName,
			data: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: testFS}
			},
		},
		"write to file url": {
			url:  testFileUrl,
			data: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: testFS}
			},
		},
		"mkdir error": {
			url:         testFileName,
			expectedErr: fmt.Errorf("couldn't create directories for file %s: operation not permitted", testFileName),
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: afero.NewReadOnlyFs(testFS)}
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fsUtils := tc.setupFileSystem()
			si := StorageInstance{fsUtils}

			err := si.WriteFromBytes(tc.url, []byte(tc.data))

			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
				writtenData, _ := si.ReadAsString(tc.url)
				require.Equal(t, writtenData, tc.data)
			}
		})
	}
}

func TestStorageInstance_WriteFromString(t *testing.T) {
	testCases := map[string]struct {
		url             string
		data            string
		setupFileSystem func() *afero.Afero
		expectedErr     error
	}{
		"write to filename": {
			url:  testFileName,
			data: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: testFS}
			},
		},
		"write to file url": {
			url:  testFileUrl,
			data: testData,
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: testFS}
			},
		},
		"mkdir error": {
			url:         testFileName,
			expectedErr: fmt.Errorf("couldn't create directories for file %s: operation not permitted", testFileName),
			setupFileSystem: func() *afero.Afero {
				testFS := afero.NewMemMapFs()
				return &afero.Afero{Fs: afero.NewReadOnlyFs(testFS)}
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fsUtils := tc.setupFileSystem()
			si := StorageInstance{fsUtils}

			err := si.WriteFromString(tc.url, tc.data)

			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
				writtenData, _ := si.ReadAsString(tc.url)
				require.Equal(t, writtenData, tc.data)
			}
		})
	}
}
