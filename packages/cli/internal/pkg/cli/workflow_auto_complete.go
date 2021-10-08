package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/spf13/cobra"
)

type WorkflowAutoComplete struct {
	workflowManagerFactory func() workflow.Interface
}

func NewWorkflowAutoComplete() *WorkflowAutoComplete {
	return &WorkflowAutoComplete{
		workflowManagerFactory: func() workflow.Interface {
			return workflow.NewManager(profile)
		},
	}
}

func (w *WorkflowAutoComplete) GetWorkflowAutoComplete() func(command *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		workflows, err := w.workflowManagerFactory().ListWorkflows()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		workflowNames := make([]string, 0)
		for k := range workflows {
			workflowNames = append(workflowNames, k)
		}
		return workflowNames, cobra.ShellCompDirectiveNoFileComp
	}
}
