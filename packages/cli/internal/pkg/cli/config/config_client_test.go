package config

import (
	"errors"
	"testing"

	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	expectedError := errors.New("some error")
	mockOs.EXPECT().UserHomeDir().Return("", expectedError)
	_, err := DetermineHomeDir()

	assert.Error(t, err, expectedError.Error())
}
