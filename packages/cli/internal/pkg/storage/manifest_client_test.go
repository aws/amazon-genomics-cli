package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/stretchr/testify/assert"
)

func TestGetManifestFilepath(t *testing.T) {
	expectedManifestFilePath := "directory/path/MANIFEST.json"
	actualManifestFilePath := getManifestFilepath("directory/path")
	assert.Equal(t, expectedManifestFilePath, actualManifestFilePath)
}

func TestDoesManifestExistInDirectory(t *testing.T) {
	backupOsStat := osStat
	defer func() {
		osStat = backupOsStat
	}()

	tests := map[string]struct {
		setupMocks     func()
		input          string
		expectedOutput bool
	}{
		"exists in directory": {
			setupMocks: func() {
				osStat = func(directory string) (os.FileInfo, error) {
					if directory != "my/directory/MANIFEST.json" {
						return nil, errors.New("jsonError")
					}
					return nil, nil
				}
			},
			input:          "my/directory",
			expectedOutput: true,
		},
		"does not exist in directory": {
			setupMocks: func() {
				osStat = func(directory string) (os.FileInfo, error) {
					return nil, errors.New("jsonError")
				}
			},
			input:          "bad directory",
			expectedOutput: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMocks()
			actualOutput := DoesManifestExistInDirectory(tt.input)
			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}

func TestReadManifestInDirectory(t *testing.T) {
	backupSpecFromJson := specFromJson
	defer func() {
		specFromJson = backupSpecFromJson
	}()

	tests := map[string]struct {
		setupMocks func()
		input      string
		errMessage string
	}{
		"Read Manifest success": {
			setupMocks: func() {
				specFromJson = func(directory string) (spec.Manifest, error) {
					if directory != "my/directory/MANIFEST.json" {
						return spec.Manifest{}, errors.New("jsonError")
					}
					return spec.Manifest{}, nil
				}
			},
			input: "my/directory",
		},
		"Read Manifest failure": {
			setupMocks: func() {
				specFromJson = func(directory string) (spec.Manifest, error) {
					if directory != "bad directory/MANIFEST.json" {
						return spec.Manifest{}, errors.New("bad match")
					}
					return spec.Manifest{}, errors.New("read error")
				}
			},
			input:      "bad directory",
			errMessage: "read error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMocks()
			_, err := ReadManifestInDirectory(tt.input)
			if err != nil {
				assert.EqualError(t, err, tt.errMessage)
			}
		})
	}
}
