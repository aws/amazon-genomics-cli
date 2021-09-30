package cfn

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type DeletionResult struct {
	Error error
}

type stackDeletionTracker chan DeletionResult

var sleepDuration = time.Second * 5

func (c Client) DeleteStack(stackId string) (chan DeletionResult, error) {
	_, err := c.cfn.DeleteStack(context.Background(), &cloudformation.DeleteStackInput{
		StackName: aws.String(stackId),
	})
	if err != nil {
		return nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}

	tracker := c.newStackDeletionTracker(stackId)

	return tracker, err
}

func (c Client) newStackDeletionTracker(stackId string) stackDeletionTracker {
	tracker := make(stackDeletionTracker)
	go func() {
		defer func() { close(tracker) }()
		for {
			time.Sleep(sleepDuration)
			done := tracker.evalDeletionStatus(c.GetStackStatus(stackId))
			if done {
				return
			}
		}
	}()

	return tracker
}

func (tracker stackDeletionTracker) evalDeletionStatus(status types.StackStatus, err error) bool {
	if err != nil {
		tracker <- DeletionResult{Error: err}
		return true
	}
	switch status {
	case types.StackStatusDeleteInProgress:
		return false
	case types.StackStatusDeleteComplete:
		tracker <- DeletionResult{}
		return true
	case types.StackStatusDeleteFailed:
		tracker <- DeletionResult{Error: fmt.Errorf("failed to delete stack")}
		return true
	default:
		tracker <- DeletionResult{Error: fmt.Errorf("unexpected status of the stack: %s", status)}
		return true
	}
}
