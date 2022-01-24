package osutils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
)

var osUserHomeDir = os.UserHomeDir
var osOpen = os.Open
var osCreate = os.Create
var ioCopy = io.Copy
var filepathWalkDir = filepath.WalkDir

// DetermineHomeDir returns the file system directory where the AGC files live.
func DetermineHomeDir() (string, error) {
	dir, err := osUserHomeDir()
	if err != nil {
		return "", actionableerror.New(err, "Please check that your home or user profile directory is defined within your environment variables")
	}
	return dir, nil
}

// ExpandHomeDir returns the expanded home directory path for the current user
func ExpandHomeDir(rootPath string) (string, error) {
	if rootPath == "~" {
		homeDir, err := DetermineHomeDir()
		if err != nil {
			return "", err
		}
		rootPath = homeDir
	} else if strings.HasPrefix(rootPath, "~/") {
		homeDir, err := DetermineHomeDir()
		if err != nil {
			return "", err
		}
		rootPath = filepath.Join(homeDir, rootPath[2:])
		return rootPath, nil
	}
	return rootPath, nil
}

func EnsureDirExistence(dirPath string) error {
	dirStat, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0744)
		return err
	}

	if !dirStat.IsDir() {
		return fmt.Errorf("'%s' should be a directory", dirPath)
	}

	return err
}

func CopyFileRecursivelyToLocation(destinationDir string, sourceDir string) error {
	err := filepathWalkDir(sourceDir, func(currentPath string, dirEntry fs.DirEntry, err error) error {
		if dirEntry == nil {
			// There are several use cases when it can happen:
			// 1. provided path doesn't exist
			// 2. file or sub-directory got deleted after being listed by WalkDir
			return fmt.Errorf("file '%s' doesn't exist", currentPath)
		}
		if !dirEntry.IsDir() {
			srcFile, err := osOpen(currentPath)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			relativePath := fmt.Sprintf("%s/%s", destinationDir, dirEntry.Name())
			dstFile, err := osCreate(relativePath)
			if err != nil {
				return err
			}
			_, err = ioCopy(dstFile, srcFile)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
