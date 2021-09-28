package cdk

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
	"github.com/rs/zerolog/log"
)

func (client Client) ClearContext(appDir string) error {
	log.Debug().Msgf("executeCDKClearContext(%s)", appDir)

	cmdArgs := []string{
		"context",
		"--clear",
	}

	stream, err := ExecuteCdkCommand(appDir, cmdArgs)
	if err != nil {
		return actionable.FindSuggestionForError(err, actionable.AwsErrorMessageToSuggestedActionMap)
	}
	for event := range stream {
		if event.Err != nil {
			return actionable.FindSuggestionForError(event.Err, actionable.AwsErrorMessageToSuggestedActionMap)
		}
	}

	return nil
}
