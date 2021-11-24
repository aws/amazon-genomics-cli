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

func TestConfig_UserId(t *testing.T) {

	testCases := []struct {
		name         string
		emailAddress string
		userId       string
	}{
		{
			name:         "happy case",
			emailAddress: "user@example.com",
			userId:       "user43G9Hd",
		},
		{
			name:         "lower casing the email",
			emailAddress: "USER@EXAMPLE.COM",
			userId:       "user43G9Hd", // same as for lowercase
		},
		{
			name:         "sanitizing non alpha num",
			emailAddress: "u-se.r@example.com",
			userId:       "user3cp566",
		},
		{
			name:         "unicode chars in email",
			emailAddress: "USE😃R@EXAMPLE.COM",
			userId:       "userRx00L",
		},
		{
			name:         "cutting username at 10 chars",
			emailAddress: "userWithPrettyLongNameInEmailAddress@EXAMPLE.COM",
			userId:       "userwithpr4n50vD",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			emailAddress := testCase.emailAddress
			expectedUserId := testCase.userId
			actualUserId := userIdFromEmailAddress(emailAddress)

			assert.Equal(t, expectedUserId, actualUserId)
		})
	}
}

func TestGetFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileReader := iomocks.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().ReadFile(testFileName).Return([]byte(expectedConfigYaml), nil)

	origReadFile := readFile
	readFile = mockFileReader.ReadFile
	defer func() { readFile = origReadFile }()
	var client = Client{
		configFilePath: testFileName,
	}
	configFormat, err := client.GetFormat()
	require.NoError(t, err)
	assert.Equal(t, expectedConfig.Format.Name, configFormat)
}

func TestSetFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileReader := iomocks.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().ReadFile(testFileName).Return([]byte(expectedConfigYaml), nil)

	origReadFile := readFile
	readFile = mockFileReader.ReadFile
	defer func() { readFile = origReadFile }()

	expectedConfigBytes, _ := yaml.Marshal(expectedConfig)
	mockFileWriter := iomocks.NewMockFileWriter(ctrl)
	mockFileWriter.EXPECT().WriteFile(testFileName, expectedConfigBytes, fs.FileMode(0644)).Return(nil)

	origWriteFile := writeFile
	writeFile = mockFileWriter.WriteFile
	defer func() { writeFile = origWriteFile }()

	var client = Client{
		configFilePath: testFileName,
	}
	err := client.SetFormat(defaultFormat)
	require.NoError(t, err)
}
