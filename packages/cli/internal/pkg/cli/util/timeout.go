package util

import (
	"errors"
	"time"
)

const (
	timeoutError = "the deployment is taking longer than expected. Please review the stack in CloudFormation"
)

func DeployWithTimeout(timeoutFunction func() error, timeoutDuration time.Duration) error {
	// channel to mark when a deployment successfully completes
	completionChannel := make(chan error)
	go func() {
		err := timeoutFunction()
		completionChannel <- err
	}()
	select {
	case err := <-completionChannel:
		return err
	case <-time.After(timeoutDuration):
		return errors.New(timeoutError)
	}
}
