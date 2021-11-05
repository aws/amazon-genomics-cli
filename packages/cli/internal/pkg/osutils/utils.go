package osutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
)

var osUserHomeDir = os.UserHomeDir

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
		homeDir, _ := DetermineHomeDir()
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
