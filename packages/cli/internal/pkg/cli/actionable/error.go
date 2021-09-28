package actionable

import (
	"fmt"
	"strings"
)

type Error struct {
	Cause           error
	SuggestedAction string
}

func (e *Error) Error() string {
	return fmt.Sprintf("an error occurred caused by: %s\nsuggestion: %s\n", e.Cause, e.SuggestedAction)
}

func NewError(cause error, suggestedAction string) error {
	err := new(Error)
	err.SuggestedAction = suggestedAction
	err.Cause = cause

	return err
}

func FindSuggestionForError(cause error, errorToSuggestionMap map[string]string) error {
	if cause == nil {
		return nil
	}

	err := new(Error)
	err.Cause = cause

	errorMessage := cause.Error()
	for expectedErrorMessage, suggestedAction := range errorToSuggestionMap {
		if strings.Contains(errorMessage, expectedErrorMessage) {
			err.SuggestedAction = suggestedAction
			return err
		}
	}

	return cause
}
