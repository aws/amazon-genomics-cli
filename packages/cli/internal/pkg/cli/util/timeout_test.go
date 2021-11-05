package util

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"
)

var (
	testError = errors.New("TestError")
)

func testFunc() error {
	time.Sleep(30 * time.Millisecond)
	return testError
}

func createSuccessful() (types.StackStatus, error) {
	return types.StackStatusCreateComplete, nil
}

func createFailed() (types.StackStatus, error) {
	return types.StackStatusCreateFailed, nil
}

func TestTimeout_DeployWithTimeout_NoTimeoutError(t *testing.T) {
	err := DeployWithTimeout(testFunc, createSuccessful, 50*time.Millisecond)
	assert.EqualError(t, err, testError.Error())
}

func TestTimeout_DeployWithTimeout_TimeoutExceededNoError(t *testing.T) {
	err := DeployWithTimeout(testFunc, createSuccessful, 10*time.Millisecond)
	assert.NoError(t, err)
}

func TestTimeout_DeployWithTimeout_TimeoutExceededError(t *testing.T) {
	err := DeployWithTimeout(testFunc, createFailed, 10*time.Millisecond)
	assert.EqualError(t, err, fmt.Sprintf(timeoutError, types.StackStatusCreateFailed))
}
