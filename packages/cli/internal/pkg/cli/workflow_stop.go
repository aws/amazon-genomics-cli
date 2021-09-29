// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type stopWorkflowVars struct {
	WorkflowInstanceId string
}

type stopWorkflowOpts struct {
	stopWorkflowVars
	wfManager *workflow.Manager
}

func newStopWorkflowOpts(vars stopWorkflowVars) (*stopWorkflowOpts, error) {
	return &stopWorkflowOpts{
		stopWorkflowVars: vars,
		wfManager:        workflow.NewManager(profile),
	}, nil
}

func (o *stopWorkflowOpts) Validate() error {
	return nil
}

func (o *stopWorkflowOpts) Execute() {
	o.wfManager.StopWorkflowInstance(o.WorkflowInstanceId)
}

func BuildWorkflowStopCommand() *cobra.Command {
	vars := stopWorkflowVars{}
	cmd := &cobra.Command{
		Use:   "stop workflow_instance_id",
		Short: "Stop the workflow with the specified workflow instance id.",
		Long: `
Stop the workflow with the specified workflow instance id. Signals to the workflow engine that all running tasks of the
workflow instance should be stopped and any pending tasks should be cancelled.`,
		Example: `
agc workflow stop ae12347654329`,
		Args: cobra.ExactArgs(1),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			vars.WorkflowInstanceId = args[0]
			opts, err := newStopWorkflowOpts(vars)
			if err != nil {
				return err
			}
			log.Info().Msgf("Stopping workflow. Workflow instance id: '%s'", opts.WorkflowInstanceId)
			if err := opts.Validate(); err != nil {
				return err
			}
			opts.Execute()
			return nil
		}),
	}
	return cmd
}
