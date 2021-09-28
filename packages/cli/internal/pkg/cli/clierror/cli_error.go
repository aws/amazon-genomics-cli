package clierror

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
	"github.com/rs/zerolog/log"
)

type Error struct {
	Command         string
	CommandVars     interface{}
	Cause           error
	SuggestedAction string
}

func (e *Error) Error() string {
	if e.SuggestedAction != "" {
		return fmt.Sprintf("an error occurred invoking '%s'\nwith variables: %+v\ncaused by: %s\nsuggestion: %s",
			e.Command, e.CommandVars, e.Cause, e.SuggestedAction)

	}
	return fmt.Sprintf("an error occurred invoking '%s'\nwith variables: %+v\ncaused by: %s",
		e.Command, e.CommandVars, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// New constructs an error message intended to be read by the CLI user. Holds context in the form of the invoked command (command),
// the variables of the command (commandVars), the causal error (cause), and any suggested action the user might take to correct the problem (suggestedAction).
func New(command string, commandVars interface{}, err error) *Error {
	actionableError, ok := err.(*actionable.Error)
	if ok {
		log.Err(actionableError.Cause).Send()

		return &Error{
			Command:         command,
			CommandVars:     commandVars,
			Cause:           actionableError.Cause,
			SuggestedAction: actionableError.SuggestedAction,
		}
	}

	log.Err(err).Send()
	return &Error{
		Command:     command,
		CommandVars: commandVars,
		Cause:       err,
	}
}
