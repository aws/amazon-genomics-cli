// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribeContextOpts_Execute(t *testing.T) {
	testCases := map[string]struct {
		contextName string
		expected    types.Context
		expectedErr error
		setupMocks  func(opts *describeContextOpts)
	}{
		"valid context name": {
			contextName: testContextName1,
			expected: types.Context{
				Name:   testContextName1,
				Status: "STARTED",
				Output: types.OutputLocation{Url: "s3://some-bucket/project/TestProject/context/test-context-name-1"},
			},
			setupMocks: func(opts *describeContextOpts) {
				opts.ctxManager.(*contextmocks.MockContextManager).EXPECT().Info(testContextName1).Return(context.Detail{
					Summary:        context.Summary{Name: testContextName1},
					Status:         context.StatusStarted,
					BucketLocation: "s3://some-bucket/project/TestProject/context/test-context-name-1",
				}, nil)
			},
		},
		"info error": {
			contextName: testContextName1,
			expectedErr: fmt.Errorf("some info error"),
			setupMocks: func(opts *describeContextOpts) {
				opts.ctxManager.(*contextmocks.MockContextManager).EXPECT().Info(testContextName1).Return(context.Detail{
					Summary: context.Summary{Name: testContextName1},
					Status:  context.StatusNotStarted,
				}, fmt.Errorf("some info error"))
			},
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCtxManager := contextmocks.NewMockContextManager(ctrl)
			opts := &describeContextOpts{
				ctxManager:          mockCtxManager,
				describeContextVars: describeContextVars{tt.contextName},
			}

			tt.setupMocks(opts)
			actual, err := opts.Execute()

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestDescribeContextOpts_Validate(t *testing.T) {
	testCases := map[string]struct {
		contextName     string
		expectedContext string
		contextNameArgs []string
		expectedErr     error
	}{
		"valid single context name": {
			contextName:     testContextName1,
			expectedContext: testContextName1,
			expectedErr:     nil,
		},
		"valid single context arg": {
			contextNameArgs: []string{testContextName2},
			expectedContext: testContextName2,
			expectedErr:     nil,
		},
		"both context and context arg throws error": {
			contextName:     testContextName1,
			contextNameArgs: []string{testContextName2},
			expectedErr:     fmt.Errorf("either the '-c' flag or a context must be provided, but not both"),
		},
		"no contexts supplied": {
			expectedErr: fmt.Errorf("a context must be provided"),
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCtxManager := contextmocks.NewMockContextManager(ctrl)
			opts := &describeContextOpts{
				ctxManager:          mockCtxManager,
				describeContextVars: describeContextVars{tt.contextName},
			}

			err := opts.Validate(tt.contextNameArgs)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedContext, opts.ContextName)
			}
		})
	}
}
