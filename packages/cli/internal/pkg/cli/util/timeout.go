package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

var (
	timeoutError = "the deployment is taking longer than expected. Current deployment status is %s. Please review the stack in CloudFormation"
)

func DeployWithTimeout(timeoutFunction func() error, deploymentStackStatus func() (types.StackStatus, error), timeoutDuration time.Duration) error {
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
		status, err := deploymentStackStatus()
		if err != nil {
			return err
		}
		if status == types.StackStatusCreateComplete {
			return nil
		}
		return errors.New(fmt.Sprintf(timeoutError, string(status)))
	}
}
