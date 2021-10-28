package cli

import (
	"fmt"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testContextName = "test-context-name"
)

func TestListContextOpts_Execute(t *testing.T) {
	tests := map[string]struct {
		setExpectations func(manager *contextmocks.MockContextManager)
		expected        []types.ContextSummary
		expectedErr     error
	}{
		"no contexts": {
			setExpectations: func(ctxManager *contextmocks.MockContextManager) {
				ctxManager.EXPECT().List().Return(map[string]context.Summary{}, nil)
			},
			expected: nil,
		},
		"initial context": {
			setExpectations: func(ctxManager *contextmocks.MockContextManager) {
				ctxManager.EXPECT().List().Return(map[string]context.Summary{
					testContextName: {Name: testContextName, Engines: []spec.Engine{{Type: "type", Engine: "engine"}}},
				}, nil)
			},
			expected: []types.ContextSummary{{Name: testContextName, EngineName: "engine"}},
		},
		"list error": {
			setExpectations: func(ctxManager *contextmocks.MockContextManager) {
				ctxManager.EXPECT().List().Return(nil, fmt.Errorf("some list error"))
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
			opts := &listContextOpts{
				ctxManager:      mockContextManager,
				listContextVars: listContextVars{},
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
