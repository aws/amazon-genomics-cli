package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testError = errors.New("TestError")
)

func testFunc() error {
	return testError
}

func TestTimeout_DeployWithTimeout(t *testing.T) {
	err := DeployWithTimeout(testFunc)
	assert.EqualError(t, err, testError.Error())
}
