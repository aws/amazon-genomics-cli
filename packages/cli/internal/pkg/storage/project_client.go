package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
)

const ProjectSpecFileName = "agc-project.yaml"

type FSProjectClient struct {
	RootPath string
}

func NewProjectClientInCurrentDir() (*FSProjectClient, error) {
	projectDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return NewProjectClientWithLocation(projectDirectory)
}

func NewProjectClient() (*FSProjectClient, error) {
	projectDirectory, err := findProjectDirectoryUpwards()
	if err != nil {
		return nil, err
	}
	return NewProjectClientWithLocation(projectDirectory)
}

func findProjectDirectoryUpwards() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for currentDir != "" {
		if hasProjectFile(currentDir) {
			return currentDir, nil
		}

		currentDir = getParentDir(currentDir)
	}

	return "", actionableerror.New(fmt.Errorf("can not find '%s' file", ProjectSpecFileName), "change to a project directory or create a project spec with 'agc project init'")
}

func hasProjectFile(dirPath string) bool {
	potentialSpecPath := dirToSpecPath(dirPath)

	fileStat, err := os.Stat(potentialSpecPath)
	if os.IsNotExist(err) {
		return false
	}

	return !fileStat.IsDir()
}

func getParentDir(dirPath string) string {
	if dirPath == "" || dirPath == "/" {
		return ""
	}
	return filepath.Dir(dirPath)
}

func NewProjectClientWithLocation(dirPath string) (*FSProjectClient, error) {
	if err := validatePath(dirPath); err != nil {
		return nil, err
	}
	return &FSProjectClient{dirPath}, nil
}

func validatePath(dirPath string) error {
	dirStat, err := os.Stat(dirPath)
	if err != nil {
		return err
	}
	if !dirStat.IsDir() {
		return fmt.Errorf("'%s' should be an existing directory", dirPath)
	}
	return nil
}

func (c FSProjectClient) Read() (spec.Project, error) {
	return spec.FromYaml(dirToSpecPath(c.RootPath))
}

func (c FSProjectClient) Write(projectSpec spec.Project) error {
	return spec.ToYaml(dirToSpecPath(c.RootPath), projectSpec)
}

func (c FSProjectClient) IsInitialized() (bool, error) {
	fileName := dirToSpecPath(c.RootPath)
	specStat, err := os.Stat(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if specStat.IsDir() {
		return false, fmt.Errorf("'%s' should be a file", fileName)
	}
	return true, nil
}

func (c FSProjectClient) GetProjectName() (string, error) {
	projectSpec, err := spec.FromYaml(dirToSpecPath(c.RootPath))
	if err != nil {
		return "", err
	}
	return projectSpec.Name, nil
}

func (c FSProjectClient) GetLocation() string {
	return c.RootPath
}

func dirToSpecPath(dirPath string) string {
	return filepath.Join(dirPath, ProjectSpecFileName)
}
