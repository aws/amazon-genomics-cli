package cdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClearContext(t *testing.T) {
	client := NewClient(testDeployProfile)
	realExecuteCdkCommand := ExecuteCdkCommand
	ExecuteCdkCommand = func(appDir string, cmdArgs []string, uniqueKey string) (ProgressStream, error) {
		stream := make(ProgressStream)
		assert.Equal(t, testDeployAppPath, appDir)
		assert.Equal(t, []string{
			"context",
			"--clear",
		}, cmdArgs)
		close(stream)
		return stream, nil
	}
	defer func() { ExecuteCdkCommand = realExecuteCdkCommand }()

	err := client.ClearContext(testDeployAppPath)
	assert.NoError(t, err)
}
