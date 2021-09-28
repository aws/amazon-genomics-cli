package clierror

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
	"github.com/stretchr/testify/assert"
)

func Test_New_ActionableError(t *testing.T) {
	command, commandVars, cause, suggestion := "agc deploy context", "-c context", errors.New("some error"), "my suggestion"
	actionableError := actionable.NewError(cause, suggestion)
	actualCliError := New("agc deploy context", "-c context", actionableError)

	expectedCliError := &Error{
		Command:         command,
		CommandVars:     commandVars,
		Cause:           cause,
		SuggestedAction: suggestion,
	}
	assert.Equal(t, expectedCliError, actualCliError)
}

func Test_New_Error(t *testing.T) {
	command, commandVars, cause := "agc deploy context", "-c context", errors.New("some error")
	actualCliError := New("agc deploy context", "-c context", cause)

	expectedCliError := &Error{
		Command:     command,
		CommandVars: commandVars,
		Cause:       cause,
	}
	assert.Equal(t, expectedCliError, actualCliError)
}
