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

var testContextInfoStruct1 = context.Detail{Summary: context.Summary{Name: testContextName1}}
var testContextInfoStruct2 = context.Detail{Summary: context.Summary{Name: testContextName2}}

func TestDeployContextOpts_Validate_ValidContexts(t *testing.T) {
	opts := &deployContextOpts{deployContextVars: deployContextVars{contexts: []string{testContextName1}}}
	assert.NoError(t, opts.Validate())
}

func TestDeployContextOpts_Validate_ValidAll(t *testing.T) {
	opts := &deployContextOpts{deployContextVars: deployContextVars{deployAll: true}}
	assert.NoError(t, opts.Validate())
}

func TestDeployContextOpts_Validate_InvalidNone(t *testing.T) {
	opts := &deployContextOpts{deployContextVars: deployContextVars{}}
	assert.Error(t, opts.Validate())
}

func TestDeployContextOpts_Validate_InvalidBoth(t *testing.T) {
	opts := &deployContextOpts{deployContextVars: deployContextVars{deployAll: true, contexts: []string{testContextName1}}}
	assert.Error(t, opts.Validate())
}

func TestDeployContextOpts_Execute_One(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().Deploy(testContextName1, true).Return(nil)
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
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	ctxMock.EXPECT().Deploy(testContextName1, true).Return(nil)
	ctxMock.EXPECT().Deploy(testContextName2, true).Return(nil)
	ctxMock.EXPECT().Info(testContextName1).Return(testContextInfoStruct1, nil)
	ctxMock.EXPECT().Info(testContextName2).Return(testContextInfoStruct2, nil)

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	info, err := opts.Execute()
	require.NoError(t, err)
	require.Equal(t, []context.Detail{testContextInfoStruct1, testContextInfoStruct2}, info)
}

func TestDeployContextOpts_Execute_ListError(t *testing.T) {
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
	_, err := opts.Execute()
	require.Equal(t, expectedErr, err)
}

func TestDeployContextOpts_Execute_InfoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	expectedErr := errors.New("one or more contexts failed to deploy")
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	ctxMock.EXPECT().Deploy(testContextName1, true).Return(nil)
	ctxMock.EXPECT().Deploy(testContextName2, true).Return(nil)
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{}, errors.New("some info error"))
	ctxMock.EXPECT().Info(testContextName2).Return(testContextInfoStruct2, nil)

	opts := &deployContextOpts{
		deployContextVars: deployContextVars{deployAll: true},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	_, err := opts.Execute()
	require.Equal(t, expectedErr, err)
}
