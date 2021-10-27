// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testContextInfoStruct1 = context.Detail{Summary: context.Summary{Name: testContextName1}}
var testContextInfoStruct2 = context.Detail{Summary: context.Summary{Name: testContextName2}}

func TestDeployContextOpts_Validate_ValidContexts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	opts := &deployContextOpts{
		deployContextVars: deployContextVars{},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
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
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
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
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
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
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	err := opts.Validate([]string{})
	require.Equal(t, expectedErr, err)
}

func TestDeployContextOpts_Execute_One(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().Deploy([]string{testContextName1}).Return(nil)
	ctxMock.EXPECT().Info(testContextName1).Return(testContextInfoStruct1, nil)
	opts := &deployContextOpts{
		deployContextVars: deployContextVars{contexts: []string{testContextName1}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	info, err := opts.Execute()
	require.NoError(t, err)
	require.Equal(t, []context.Detail{testContextInfoStruct1}, info)
}

func TestDeployContextOpts_Execute_All(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().Deploy([]string{testContextName1, testContextName2}).Return([]context.ProgressResult{{Context: testContextName1}, {Context: testContextName2}})
	ctxMock.EXPECT().Info(testContextName1).Return(testContextInfoStruct1, nil)
	ctxMock.EXPECT().Info(testContextName2).Return(testContextInfoStruct2, nil)

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true, contexts: []string{testContextName1, testContextName2}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	info, err := opts.Execute()
	require.NoError(t, err)
	expectedDetailList := []context.Detail{testContextInfoStruct1, testContextInfoStruct2}
	require.Equal(t, expectedDetailList, info)
}

func TestDeployContextOpts_ExecuteAll_LogsOutErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	progressResults := []context.ProgressResult{{Context: testContextName1, Err: errors.New("some error"), Outputs: []string{"log1", "log2"}}, {Context: testContextName2}}
	ctxMock.EXPECT().Deploy([]string{testContextName1, testContextName2}).Return(progressResults)

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true, contexts: []string{testContextName1, testContextName2}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	_, err := opts.Execute()
	require.Error(t, err)
}

func TestDeployContextOpts_Execute_InfoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	expectedErr := actionableerror.New(errors.New("one or more contexts failed to deploy"), "")
	ctxMock.EXPECT().Deploy([]string{testContextName1, testContextName2}).Return(nil)
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{}, errors.New("some info error"))
	ctxMock.EXPECT().Info(testContextName2).Return(testContextInfoStruct2, nil)

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true, contexts: []string{testContextName1, testContextName2}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	_, err := opts.Execute()
	require.Equal(t, expectedErr, err)
}
