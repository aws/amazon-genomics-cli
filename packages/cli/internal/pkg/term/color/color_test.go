// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package color

import (
	"testing"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
)

type envVar struct {
	env map[string]string
}

func (e *envVar) lookupEnv(key string) (string, bool) {
	v, ok := e.env[key]
	return v, ok
}

func TestColorEnvVarSetToFalse(t *testing.T) {
	env := &envVar{
		env: map[string]string{colorEnvVar: "false"},
	}
	lookupEnv = env.lookupEnv

	DisableColorBasedOnEnvVar()

	require.True(t, core.DisableColor, "expected to be true when COLOR is disabled")
	require.True(t, color.NoColor, "expected to be true when COLOR is disabled")
}

func TestColorEnvVarSetToTrue(t *testing.T) {
	env := &envVar{
		env: map[string]string{colorEnvVar: "true"},
	}
	lookupEnv = env.lookupEnv

	DisableColorBasedOnEnvVar()

	require.False(t, core.DisableColor, "expected to be false when COLOR is enabled")
	require.False(t, color.NoColor, "expected to be true when COLOR is enabled")
}

func TestColorEnvVarNotSet(t *testing.T) {
	env := &envVar{
		env: make(map[string]string),
	}
	lookupEnv = env.lookupEnv

	DisableColorBasedOnEnvVar()

	require.Equal(t, core.DisableColor, color.NoColor, "expected to be the same as color.NoColor")
}

func TestHelp(t *testing.T) {
	output := Help("test")
	require.Equal(t, "\u001B[2mtest\u001B[0m", output)
}

func TestEmphasize(t *testing.T) {
	output := Emphasize("test")
	require.Equal(t, "\u001B[1mtest\u001B[0m", output)
}

func TestHighlightUserInput(t *testing.T) {
	output := HighlightUserInput("test")
	require.Equal(t, "\u001B[1mtest\u001B[0m", output)
}

func TestHighlightResource(t *testing.T) {
	output := HighlightResource("test")
	require.Equal(t, "\u001B[94mtest\u001B[0m", output)
}

func TestHighlightCode(t *testing.T) {
	output := HighlightCode("test")
	require.Equal(t, "\u001B[96m`test`\u001B[0m", output)
}

func TestProd(t *testing.T) {
	output := Prod("test")
	require.Equal(t, "\u001B[33;1mtest\u001B[0m", output)
}
