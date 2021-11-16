package cdk

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/rs/zerolog/log"
)

func (client Client) ClearContext(appDir string) error {
	log.Debug().Msgf("executeCDKClearContext(%s)", appDir)

	cmdArgs := []string{
		"context",
		"--clear",
	}

	stream, err := ExecuteCdkCommand(appDir, cmdArgs, "clear-context")
	if err != nil {
		return actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	for event := range stream {
		if event.Err != nil {
			return actionableerror.FindSuggestionForError(event.Err, actionableerror.AwsErrorMessageToSuggestedActionMap)
		}
	}

	return nil
}
