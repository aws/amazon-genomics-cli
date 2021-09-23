package clierror

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Error struct {
	Command         string
	CommandVars     interface{}
	Cause           error
	SuggestedAction string
}

func (e *Error) Error() string {
	return fmt.Sprintf("an error occurred invoking '%s'\nwith variables: %+v\ncaused by: %s\nsuggestion: %s",
		e.Command, e.CommandVars, e.Cause, e.SuggestedAction)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// New constructs an error message intended to be read by the CLI user. Holds context in the form of the invoked command (command),
// the variables of the command (commandVars), the causal error (cause), and any suggested action the user might take to correct the problem (suggestedAction).
func New(command string, commandVars interface{}, cause error, suggestedAction string) *Error {
	log.Err(cause).Send()

	actionableError, ok := cause.(ActionableError)
	if ok {
		cause = actionableError.Cause
		suggestedAction = actionableError.SuggestedAction
	}

	return &Error{
		Command:         command,
		CommandVars:     commandVars,
		Cause:           cause,
		SuggestedAction: suggestedAction,
	}
}
