// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type describeWorkflowVars struct {
	WorkflowName string
}

type describeWorkflowOpts struct {
	describeWorkflowVars
	wfManager *workflow.Manager
}

func newDescribeWorkflowOpts(vars describeWorkflowVars) (*describeWorkflowOpts, error) {
	return &describeWorkflowOpts{
		describeWorkflowVars: vars,
		wfManager:            workflow.NewManager(profile),
	}, nil
}

func (o *describeWorkflowOpts) Validate() error {
	return nil
}

func (o *describeWorkflowOpts) Execute() (types.Workflow, error) {
	details, err := o.wfManager.DescribeWorkflow(o.WorkflowName)
	if err != nil {
		return types.Workflow{}, err
	}

	workflow := types.Workflow{
		Name:         details.Name,
		TypeLanguage: details.TypeLanguage,
		TypeVersion:  details.TypeVersion,
		Source:       details.Source,
	}

	return workflow, nil
}

// BuildWorkflowDescribeCommand builds the command to describe the information for a workflow  in the current project.
func BuildWorkflowDescribeCommand() *cobra.Command {
	vars := describeWorkflowVars{}
	cmd := &cobra.Command{
		Use:   "describe workflow_name [--max_instances]",
		Short: "Show the information for a specific workflow in the current project",
		Long: `describe is for showing details on the specified workflow.
It includes workflow specification and list of recent instances of that workflow.
An instance is created every time we run a workflow.

` + DescribeOutput(types.Workflow{}),
		Args: cobra.ExactArgs(1),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			vars.WorkflowName = args[0]
			opts, err := newDescribeWorkflowOpts(vars)
			if err != nil {
				return err
			}
			log.Info().Msgf("Describing workflow '%s'.", opts.WorkflowName)
			if err := opts.Validate(); err != nil {
				return err
			}
			workflow, err := opts.Execute()
			if err != nil {
				return clierror.New("workflow describe", vars, err)
			}
			format.Default.Write(workflow)
			return nil
		}),
		ValidArgsFunction: NewWorkflowAutoComplete().GetWorkflowAutoComplete(),
	}
	return cmd
}
