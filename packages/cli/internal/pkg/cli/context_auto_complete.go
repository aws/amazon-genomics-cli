package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/spf13/cobra"
)

type ContextAutoComplete struct {
	ctxManagerFactory func() context.Interface
}

func NewContextAutoComplete() *ContextAutoComplete {
	return &ContextAutoComplete{
		ctxManagerFactory: func() context.Interface {
			return context.NewManager(profile)
		},
	}
}

func (c *ContextAutoComplete) GetContextAutoComplete() func(command *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		contexts, err := c.ctxManagerFactory().List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		contextNames := make([]string, 0)
		for k := range contexts {
			contextNames = append(contextNames, k)
		}
		return contextNames, cobra.ShellCompDirectiveNoFileComp
	}
}
