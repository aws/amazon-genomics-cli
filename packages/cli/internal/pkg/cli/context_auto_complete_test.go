package cli

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestContextAutoComplete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {}, testContextName2: {}}, nil)
	contextAutoComplete := &ContextAutoComplete{
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	autoCompleteFunction := contextAutoComplete.GetContextAutoComplete()
	_, compDirective := autoCompleteFunction(nil, make([]string, 0), "")
	assert.Equal(t, compDirective, cobra.ShellCompDirectiveNoFileComp)
}

func TestContextAutoComplete_WithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	ctxMock.EXPECT().List().Return(map[string]context.Summary{}, errors.New("Test Context Error"))
	contextAutoComplete := &ContextAutoComplete{
		ctxManagerFactory: func() context.Interface {
			return ctxMock
		},
	}
	autoCompleteFunction := contextAutoComplete.GetContextAutoComplete()
	actualKeys, compDirective := autoCompleteFunction(nil, make([]string, 0), "")
	assert.Equal(t, []string(nil), actualKeys)
	assert.Equal(t, compDirective, cobra.ShellCompDirectiveError)
}
