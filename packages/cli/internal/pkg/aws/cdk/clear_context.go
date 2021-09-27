package cdk

import "github.com/rs/zerolog/log"

func (client Client) ClearContext(appDir string) error {
	log.Debug().Msgf("executeCDKClearContext(%s)", appDir)

	cmdArgs := []string{
		"context",
		"--clear",
	}

	stream, err := ExecuteCdkCommand(appDir, cmdArgs)
	if err != nil {
		return err
	}
	for event := range stream {
		if event.Err != nil {
			return event.Err
		}
	}

	return nil
}
