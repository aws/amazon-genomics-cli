// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"sort"

	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/types"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type listWorkflowVars struct{}

type listWorkflowOpts struct {
	listWorkflowVars
	wfManager *workflow.Manager
}

func newListWorkflowOpts(vars listWorkflowVars) (*listWorkflowOpts, error) {
	return &listWorkflowOpts{
		listWorkflowVars: vars,
		wfManager:        workflow.NewManager(profile),
	}, nil
}

func (o *listWorkflowOpts) Validate() error {
	return nil
}

func (o *listWorkflowOpts) Execute() ([]types.WorkflowName, error) {
	workflowSummaries, err := o.wfManager.ListWorkflows()
	if err != nil {
		return nil, err
	}
	var workflowNames []types.WorkflowName
	for _, summary := range workflowSummaries {
		workflowNames = append(workflowNames, types.WorkflowName{Name: summary.Name})
	}
	sortWorkflowNames(workflowNames)
	return workflowNames, nil
}

func sortWorkflowNames(workflowNames []types.WorkflowName) {
	sort.Slice(workflowNames, func(i, j int) bool {
		return workflowNames[i].Name < workflowNames[j].Name
	})
}

func BuildWorkflowListCommand() *cobra.Command {
	vars := listWorkflowVars{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show a list of workflows related to the current project",
		Long: `list is for showing a combined list of workflows defined in the project specification
and workflow instances that were run in this AWS account.

` + DescribeOutput([]types.WorkflowName{}),
		Args: cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newListWorkflowOpts(vars)
			if err != nil {
				return err
			}
			log.Info().Msgf("Listing workflows.")
			if err := opts.Validate(); err != nil {
				return err
			}
			workflowNames, err := opts.Execute()
			if err != nil {
				return clierror.New("logs access", vars, err, "check you have valid aws credentials, check that a valid agc-project.yaml file exists in this directory or one of it's parent directories and that it defines at least one workflow")
			}
			format.Default.Write(workflowNames)
			return nil
		}),
	}
	return cmd
}
