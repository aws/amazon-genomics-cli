package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFSProjectClient_IsInitialized_SpecExists(t *testing.T) {
	tempDir := t.TempDir()
	specFile, err := os.Create(filepath.Join(tempDir, ProjectSpecFileName))
	require.NoError(t, err)
	_ = specFile.Close()
	client, err := NewProjectClientWithLocation(tempDir)
	require.NoError(t, err)
	actual, err := client.IsInitialized()
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestFSProjectClient_NewClient_DirIsFile(t *testing.T) {
	tempDir := t.TempDir()
	projDirName := filepath.Join(tempDir, "myproject")
	projDir, err := os.Create(projDirName)
	require.NoError(t, err)
	_ = projDir.Close()
	_, err = NewProjectClientWithLocation(projDirName)
	assert.Error(t, err)
}

func TestFSProjectClient_IsInitialized_SpecDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()
	client, err := NewProjectClientWithLocation(tempDir)
	require.NoError(t, err)
	actual, err := client.IsInitialized()
	require.NoError(t, err)
	assert.False(t, actual)
}

func TestFSProjectClient_IsInitialized_DirWithSpecName(t *testing.T) {
	tempDir := t.TempDir()
	err := os.Mkdir(filepath.Join(tempDir, ProjectSpecFileName), 0755)
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewProjectClientWithLocation(tempDir)
	require.NoError(t, err)
	_, err = client.IsInitialized()
	assert.Error(t, err)
}

func TestFSProjectClient_NewClient_CanFindProjectFileGoingUp(t *testing.T) {
	tempDir := t.TempDir()
	specFile, err := os.Create(filepath.Join(tempDir, ProjectSpecFileName))
	defer specFile.Close() //nolint:errcheck,staticcheck
	require.NoError(t, err)

	const childDirName = "child-dir"
	childDirPath := filepath.Join(tempDir, childDirName)
	err = os.Mkdir(childDirPath, 0755)
	require.NoError(t, err)

	wd, err := os.Getwd()
	defer os.Chdir(wd) //nolint:errcheck
	require.NoError(t, err)

	testCases := map[string]struct {
		path string
	}{
		"loads project file from project root directory": {
			path: tempDir,
		},
		"loads project file from a project sub directory": {
			path: childDirPath,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			os.Chdir(tc.path) //nolint:errcheck

			client, err := NewProjectClient()
			require.NoError(t, err)
			actual, err := client.IsInitialized()
			require.NoError(t, err)
			assert.True(t, actual)
		})
	}
}
