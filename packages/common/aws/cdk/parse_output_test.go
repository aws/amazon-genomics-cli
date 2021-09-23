package cdk

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testOutputsJson = `{"RSC":{"RSCAPI":"abc", "VERSION": 1}}`
const badFormatOutputsJson = `{"RSC":{"RSCAPI":}}`
const unsupportedFormatOutputsJson = `{"RSC":{"RSCAPI": false}}`
const testFileName = "testOp.json"

func setup(testPath string, js string) error {
	op := []byte(js)
	err := ioutil.WriteFile(filepath.Join(testPath, testFileName), op, 0644)
	return err
}

func TestParseOutput(t *testing.T) {
	testPath, err := os.UserHomeDir()
	err = setup(testPath, testOutputsJson)
	actualMap, err := ParseOutput(filepath.Join(testPath, testFileName))

	expected := map[string]string{
		"RSCAPI":  "abc",
		"VERSION": "1.000000",
	}
	require.NoError(t, err)
	assert.Equal(t, expected, actualMap)

	err = os.Remove(filepath.Join(testPath, testFileName))
	require.NoError(t, err)
}

func TestParseOutputWithBadJson(t *testing.T) {
	testPath, err := os.UserHomeDir()
	err = setup(testPath, badFormatOutputsJson)
	actualMap, err := ParseOutput(filepath.Join(testPath, testFileName))
	assert.Error(t, err)
	assert.Nil(t, actualMap)

	err = os.Remove(filepath.Join(testPath, testFileName))
	require.NoError(t, err)
}

func TestParseOutputWithUnsupportedJson(t *testing.T) {
	testPath, err := os.UserHomeDir()
	err = setup(testPath, unsupportedFormatOutputsJson)
	actualMap, err := ParseOutput(filepath.Join(testPath, testFileName))
	assert.Error(t, err)
	assert.Nil(t, actualMap)

	err = os.Remove(filepath.Join(testPath, testFileName))
	require.NoError(t, err)
}
