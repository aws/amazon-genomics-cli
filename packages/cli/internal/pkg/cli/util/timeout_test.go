package util

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testError = errors.New("TestError")
)

func testFunc() error {
	time.Sleep(30 * time.Millisecond)
	return testError
}

func TestTimeout_DeployWithTimeout_NoTimeoutError(t *testing.T) {
	err := DeployWithTimeout(testFunc, 50*time.Millisecond)
	assert.EqualError(t, err, testError.Error())
}

func TestTimeout_DeployWithTimeout_TimeoutError(t *testing.T) {
	err := DeployWithTimeout(testFunc, 10*time.Millisecond)
	assert.EqualError(t, err, timeoutError)
}
