package actionable

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_WithMatchingSuggestion(t *testing.T) {
	var errorMessageToSuggestedActionMap = map[string]string{
		"error occurred": "Please do this suggestion",
		"some error":     "Please do this other suggestion",
	}

	coreError := errors.New("my error occurred")
	actualError := FindSuggestionForError(coreError, errorMessageToSuggestedActionMap)

	expectedError := &Error{
		coreError,
		"Please do this suggestion",
	}

	assert.Equal(t, expectedError, actualError)
}

func Test_New_NoMatchingSuggestion(t *testing.T) {
	var errorMessageToSuggestedActionMap = map[string]string{
		"error occurred": "Please do this suggestion",
		"some error":     "Please do this other suggestion",
	}

	coreError := errors.New("some different error")
	actualError := FindSuggestionForError(coreError, errorMessageToSuggestedActionMap)

	assert.Equal(t, coreError, actualError)
}

func Test_Error_WithSuggestion(t *testing.T) {
	err := NewError(fmt.Errorf("some error"), "some suggestion")

	assert.Equal(t, "an error occurred caused by: some error\nsuggestion: some suggestion\n", err.Error())
}
