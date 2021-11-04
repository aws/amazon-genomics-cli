package util

import (
	"errors"
	"time"
)

const (
	contextDeploymentTimeout = 30 * time.Minute
)

func DeployWithTimeout(timeoutFunction func() error) error {
	// channel to mark when a deployment successfully completes
	completionChannel := make(chan error)
	go func() {
		err := timeoutFunction()
		completionChannel <- err
	}()
	select {
	case err := <-completionChannel:
		return err
	case <-time.After(contextDeploymentTimeout):
		return errors.New("deployment taking longer then expected. please review stack deployment in cloudformation")
	}
}
