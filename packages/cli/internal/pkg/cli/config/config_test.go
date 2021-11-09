package config

import (
	"io/fs"
	"testing"

	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	testFileName       = "config.yaml"
	expectedConfigYaml = `
user:
    email: my@email.com
format:
    format: text`
)

var (
	expectedConfig = Config{
		User{
			Email: "my@email.com",
		},
		Format{
			Format: "text",
		},
	}
)

func TestConfig_ReadData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileReader := iomocks.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().ReadFile(testFileName).Return([]byte(expectedConfigYaml), nil)

	origReadFile := readFile
	readFile = mockFileReader.ReadFile
	defer func() { readFile = origReadFile }()

	configData, err := fromYaml(testFileName)
	require.NoError(t, err)
	assert.Equal(t, expectedConfig, configData)
}

func TestConfig_WriteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedConfigBytes, _ := yaml.Marshal(expectedConfig)
	mockFileWriter := iomocks.NewMockFileWriter(ctrl)
	mockFileWriter.EXPECT().WriteFile(testFileName, expectedConfigBytes, fs.FileMode(0644)).Return(nil)

	origWriteFile := writeFile
	writeFile = mockFileWriter.WriteFile
	defer func() { writeFile = origWriteFile }()

	err := toYaml(testFileName, expectedConfig)
	require.NoError(t, err)
}
