package actionableerror

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

func New(cause error, suggestedAction string) error {
	err := new(Error)
	err.SuggestedAction = suggestedAction
	err.Cause = cause

	return err
}

func FindSuggestionForError(cause error, errorToSuggestionMap map[string]string) error {
	errorMessage := cause.Error()

	for expectedErrorMessage, suggestedAction := range errorToSuggestionMap {
		if strings.Contains(errorMessage, expectedErrorMessage) {
			return &Error{
				SuggestedAction: suggestedAction,
				Cause:           cause,
			}
		}
	}

	return cause
}
