package config

import (
	"testing"

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
