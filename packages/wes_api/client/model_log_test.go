/*
 * Workflow Execution Service
 * API version: 1.0.0
 *
 * Tests for custom Log unmarshalling.
 */

package wes_client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testQuotedExitCode       = `{"name": "Task1", "cmd": ["ls", "-lah"], "exit_code": "0"}`
	testUnquotedExitCode     = `{"name": "Task1", "cmd": ["ls", "-lah"], "exit_code": 0}`
	testUnsetExitCode        = `{"name": "Task1", "cmd": ["ls", "-lah"]}`
	testUnacceptableExitCode = `{"name": "Task1", "cmd": ["ls", "-lah"], "exit_code": true}`
)

func TestLog_UnmarshallJSON_QuotedExitCode(t *testing.T) {
	var unmarshallTo Log

	err := json.Unmarshal([]byte(testQuotedExitCode), &unmarshallTo)

	assert.NoError(t, err)
	assert.Equal(t, "Task1", unmarshallTo.Name)
	assert.Equal(t, 2, len(unmarshallTo.Cmd))
	assert.Equal(t, "ls", unmarshallTo.Cmd[0])
	assert.Equal(t, "-lah", unmarshallTo.Cmd[1])
	assert.Equal(t, "0", unmarshallTo.ExitCode)
}

func TestLog_UnmarshallJSON_UnquotedExitCode(t *testing.T) {
	var unmarshallTo Log

	err := json.Unmarshal([]byte(testUnquotedExitCode), &unmarshallTo)

	assert.NoError(t, err)
	assert.Equal(t, "Task1", unmarshallTo.Name)
	assert.Equal(t, 2, len(unmarshallTo.Cmd))
	assert.Equal(t, "ls", unmarshallTo.Cmd[0])
	assert.Equal(t, "-lah", unmarshallTo.Cmd[1])
	assert.Equal(t, "0", unmarshallTo.ExitCode)
}

func TestLog_UnmarshallJSON_UnsetExitCode(t *testing.T) {
	var unmarshallTo Log

	err := json.Unmarshal([]byte(testUnsetExitCode), &unmarshallTo)

	assert.NoError(t, err)
	assert.Equal(t, "Task1", unmarshallTo.Name)
	assert.Equal(t, 2, len(unmarshallTo.Cmd))
	assert.Equal(t, "ls", unmarshallTo.Cmd[0])
	assert.Equal(t, "-lah", unmarshallTo.Cmd[1])
	assert.Equal(t, "", unmarshallTo.ExitCode)
}

func TestLog_UnmarshallJSON_UnacceptableExitCode(t *testing.T) {
	var unmarshallTo Log

	err := json.Unmarshal([]byte(testUnacceptableExitCode), &unmarshallTo)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "Task1", unmarshallTo.Name)
	assert.Equal(t, 2, len(unmarshallTo.Cmd))
	assert.Equal(t, "ls", unmarshallTo.Cmd[0])
	assert.Equal(t, "-lah", unmarshallTo.Cmd[1])
	assert.Equal(t, "", unmarshallTo.ExitCode)
}
