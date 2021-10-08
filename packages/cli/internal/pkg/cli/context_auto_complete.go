package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/spf13/cobra"
)

func ContextAutoComplete(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	manager := context.NewManager(profile)
	workflows, err := manager.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	workflowNames := make([]string, len(workflows))
	for k := range workflows {
		workflowNames = append(workflowNames, k)
	}
	return workflowNames, cobra.ShellCompDirectiveNoFileComp
}
