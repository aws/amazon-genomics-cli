// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"testing"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextStatusOpts_Execute(t *testing.T) {
	tests := map[string]struct {
		current         bool
		setExpectations func(manager *contextmocks.MockContextManager)
		expected        []context.Instance
		expectedErr     error
	}{
		"execute command": {
			setExpectations: func(ctxManager *contextmocks.MockContextManager) {
				ctxManager.EXPECT().StatusList().Return([]context.Instance{{
					ContextName:            "context",
					ContextStatus:          context.StatusStarted,
					ContextReason:          "",
					IsDefinedInProjectFile: false,
				}}, nil)
			},
			expected: []context.Instance{{
				ContextName:            "context",
				ContextStatus:          context.StatusStarted,
				ContextReason:          "",
				IsDefinedInProjectFile: false,
			}},
		},
		"Context manager error": {
			setExpectations: func(ctxManager *contextmocks.MockContextManager) {
				ctxManager.EXPECT().StatusList().Return(nil, fmt.Errorf("some list error"))
			},
			expectedErr: fmt.Errorf("some list error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContextManager := contextmocks.NewMockContextManager(ctrl)
			tt.setExpectations(mockContextManager)
			opts := &contextStatusOpts{
				ctxManager: mockContextManager,
			}
			actual, err := opts.Execute()
			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
