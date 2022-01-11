package zipfile

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFile1Name = "TestFile1"
	testFile2Name = "TestFile2"
	testDir1Name  = "TestDir1"
	testDir2Name  = "TestDir2"
	testDir3Name  = "TestDir3"
)

func TestCompressToTmp(t *testing.T) {
	tmpDir := t.TempDir()
	if err := createTestFiles(tmpDir); err != nil {
		t.Fatal(err)
	}
	tests := map[string]struct {
		srcPath  string
		expected []string
	}{
		"single file": {
			srcPath:  filepath.Join(tmpDir, testFile1Name),
			expected: []string{testFile1Name},
		},
		"deep folder": {
			srcPath: tmpDir,
			expected: []string{
				testFile1Name,
				testFile2Name,
				filepath.Join(testDir1Name, testFile1Name),
				filepath.Join(testDir1Name, testFile2Name),
				filepath.Join(testDir1Name, testDir2Name, testFile1Name),
				filepath.Join(testDir1Name, testDir2Name, testFile2Name),
			},
		},
		"empty folder": {
			srcPath:  filepath.Join(tmpDir, testDir3Name),
			expected: []string{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			zipPath, err := CompressToTmp(tt.srcPath)
			require.NoError(t, err)
			actual, err := listZip(zipPath)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestCompressToTmp_SourceDoesNotExist(t *testing.T) {
	zipPath, err := CompressToTmp("FooBar")
	if assert.Error(t, err) {
		assert.EqualError(t, err, "file 'FooBar' doesn't exist")
	}
	assert.Empty(t, zipPath)
}

func createTestFiles(root string) error {
	content := []byte("Test!")
	if err := ioutil.WriteFile(filepath.Join(root, testFile1Name), content, 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(root, testFile2Name), content, 0600); err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(root, testDir1Name), 0700); err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(root, testDir1Name, testDir2Name), 0700); err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(root, testDir3Name), 0700); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(root, testDir1Name, testFile1Name), content, 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(root, testDir1Name, testFile2Name), content, 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(root, testDir1Name, testDir2Name, testFile1Name), content, 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(root, testDir1Name, testDir2Name, testFile2Name), content, 0600); err != nil {
		return err
	}
	return nil
}

func listZip(zipPath string) ([]string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range reader.File {
		names = append(names, f.Name)
	}
	return names, nil
}
