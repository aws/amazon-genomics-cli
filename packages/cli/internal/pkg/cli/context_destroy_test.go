// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	managermocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestroyContextOpts_Validate_ValidContexts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	wfMock := managermocks.NewMockWorkflowManager(ctrl)
	wfMock.EXPECT().StatusWorkflowByContext(testContextName1, 20).Return([]workflow.InstanceSummary{}, nil)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{contexts: []string{testContextName1}},
		wfsManager:         wfMock,
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	assert.NoError(t, opts.Validate())
}

func TestDestroyContextOpts_Validate_ValidAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	wfMock.EXPECT().StatusWorkflowByContext(testContextName1, workflowMaxInstanceDefault).Return([]workflow.InstanceSummary{}, nil)
	wfMock.EXPECT().StatusWorkflowByContext(testContextName2, workflowMaxInstanceDefault).Return([]workflow.InstanceSummary{}, nil)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{destroyAll: true},
		wfsManager:         wfMock,
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		}}
	assert.NoError(t, opts.Validate())
}

func TestDestroyContextOpts_Validate_InvalidNone(t *testing.T) {
	contextCtrl := gomock.NewController(t)
	defer contextCtrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(contextCtrl)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{},
		wfsManager:         wfMock,
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		}}
	assert.Error(t, opts.Validate())
}

func TestDestroyContextOpts_Validate_InvalidBoth(t *testing.T) {
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{destroyAll: true, contexts: []string{testContextName1}},
		wfsManager:         wfMock}
	assert.Error(t, opts.Validate())
}

func TestDestroyContextOpts_Validate_NonExistingContexts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName2: {}}, nil)

	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{contexts: []string{testContextName1}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	assert.Error(t, opts.Validate())
}

func TestDestroyContextOpts_Validate_ContextRetrievalFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{}, fmt.Errorf("some error"))

	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{contexts: []string{testContextName1}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	assert.Error(t, opts.Validate())
}

func TestDestroyContextOpts_Validate_ContainsRunningContext(t *testing.T) {
	contextCtrl := gomock.NewController(t)
	defer contextCtrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(contextCtrl)
	contexts := []string{testContextName1}
	failedSummary := []workflow.InstanceSummary{{State: "RUNNING"}}
	expectedError := fmt.Sprintf("Context '%s' contains running workflows. Please stop all workflows before destroying context.", testContextName1)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	wfMock.EXPECT().StatusWorkflowByContext(testContextName1, workflowMaxInstanceDefault).Return(failedSummary, nil)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{contexts: contexts},
		wfsManager:         wfMock,
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		}}
	err := opts.Validate()
	assert.Equal(t, expectedError, err.Error())
}

func TestDestroyContextOpts_GetContexts_DontGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}}, nil)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{
			destroyAll: false,
			contexts:   []string{testContextName1}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
		wfsManager: wfMock,
	}
	require.NoError(t, opts.getContexts())
}

func TestDestroyContextOpts_GetContexts_ListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	contextListMap := map[string]context.Summary{testContextName1: {}, testContextName2: {}}
	ctxMock.EXPECT().List().Return(contextListMap, nil)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{destroyAll: true},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
		wfsManager: wfMock,
	}
	require.NoError(t, opts.getContexts())
}

func TestDestroyContextOpts_GetContexts_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	expectedErr := errors.New("some list error")
	ctxMock.EXPECT().List().Return(nil, expectedErr)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{destroyAll: true},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
		wfsManager: wfMock,
	}
	err := opts.getContexts()
	require.Equal(t, expectedErr, err)
}

func TestDestroyContextOpts_Execute_One(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().Destroy(testContextName1, true).Return(nil)
	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{contexts: []string{testContextName1}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
		wfsManager: wfMock,
	}
	err := opts.Execute()
	require.NoError(t, err)
}

func TestDestroyContextOpts_Execute_All(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().Destroy(testContextName1, true).Return(nil)
	ctxMock.EXPECT().Destroy(testContextName2, true).Return(nil)

	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{destroyAll: true, contexts: []string{testContextName1, testContextName2}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
		wfsManager: wfMock,
	}
	err := opts.Execute()
	require.NoError(t, err)
}

func TestDestroyContextOpts_Execute_DestroyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	workflowCtrl := gomock.NewController(t)
	defer workflowCtrl.Finish()
	wfMock := managermocks.NewMockWorkflowManager(workflowCtrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	expectedErr := errors.New("one or more contexts failed to be destroyed")
	ctxMock.EXPECT().Destroy(testContextName1, true).Return(errors.New("some destroy error"))
	ctxMock.EXPECT().Destroy(testContextName2, true).Return(nil)

	opts := &destroyContextOpts{
		destroyContextVars: destroyContextVars{destroyAll: true, contexts: []string{testContextName1, testContextName2}},
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
		wfsManager: wfMock,
	}
	err := opts.Execute()
	require.Equal(t, expectedErr, err)
}
