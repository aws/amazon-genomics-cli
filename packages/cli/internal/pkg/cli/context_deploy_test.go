// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployContextOpts_Validate_ValidContexts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	opts := &deployContextOpts{
		deployContextVars: deployContextVars{},
		ctxManager:        ctxMock,
	}
	assert.NoError(t, opts.Validate([]string{testContextName1}))
}
func TestDeployContextOpts_Validate_ValidContexts_DeprecatedArgs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	opts := &deployContextOpts{
		deployContextVars: deployContextVars{contexts: []string{testContextName1, testContextName2}},
		ctxManager:        ctxMock,
	}
	assert.NoError(t, opts.Validate([]string{}))
}

func TestDeployContextOpts_Validate_ValidAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true},
		ctxManager:        ctxMock,
	}
	assert.NoError(t, opts.Validate([]string{}))
}

func TestDeployContextOpts_Validate_InvalidNone(t *testing.T) {
	opts := &deployContextOpts{deployContextVars: deployContextVars{}}
	assert.Error(t, opts.Validate([]string{}))
}

func TestDeployContextOpts_Validate_InvalidBoth(t *testing.T) {
	opts := &deployContextOpts{deployContextVars: deployContextVars{deployAll: true}}
	assert.Error(t, opts.Validate([]string{testContextName1}))
}

func TestDeployContextOpts_Validate_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	expectedErr := errors.New("some list error")
	ctxMock.EXPECT().List().Return(nil, expectedErr)

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true},
		ctxManager:        ctxMock,
	}
	err := opts.Validate([]string{})
	require.Equal(t, expectedErr, err)
}

func TestDeployContextOpts_Execute_One(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().Deploy([]string{testContextName1}).Return(nil)
	opts := &deployContextOpts{
		deployContextVars: deployContextVars{contexts: []string{testContextName1}},
		ctxManager:        ctxMock,
	}
	err := opts.Execute()
	require.NoError(t, err)
}

func TestDeployContextOpts_Execute_All(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().Deploy([]string{testContextName1, testContextName2}).Return([]context.ProgressResult{{Context: testContextName1}, {Context: testContextName2}})

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true, contexts: []string{testContextName1, testContextName2}},
		ctxManager:        ctxMock,
	}
	err := opts.Execute()
	require.NoError(t, err)
}

func TestDeployContextOpts_ExecuteAll_LogsOutErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	progressResults := []context.ProgressResult{{Context: testContextName1, Err: errors.New("some error"), Outputs: []string{"log1", "log2"}}, {Context: testContextName2}}
	ctxMock.EXPECT().Deploy([]string{testContextName1, testContextName2}).Return(progressResults)

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true, contexts: []string{testContextName1, testContextName2}},
		ctxManager:        ctxMock,
	}
	err := opts.Execute()
	require.Error(t, err)
}
