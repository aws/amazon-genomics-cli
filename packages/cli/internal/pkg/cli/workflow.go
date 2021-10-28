// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/group"
	"github.com/spf13/cobra"
)

// BuildWorkflowCommand builds the top level workflow command and related subcommands.
func BuildWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "workflow",
		Short: `Commands for workflows.
Workflows are potentially-dynamic graphs of computational tasks to execute.`,
		Long: `Commands for workflows.
Workflows are potentially-dynamic graphs of computational tasks to execute.

Workflow specifications are files whose content specify a workflow to execute 
given a particular set of input parameters to use. Workflow specifications are 
typed according to which workflow definition language they use (e.g. WDL).
`,
	}

	cmd.AddCommand(BuildWorkflowRunCommand())
	cmd.AddCommand(BuildWorkflowListCommand())
	cmd.AddCommand(BuildWorkflowStatusCommand())
	cmd.AddCommand(BuildWorkflowDescribeCommand())
	cmd.AddCommand(BuildWorkflowStopCommand())
	cmd.AddCommand(BuildWorkflowOutputCommand())

	cmd.SetUsageTemplate(template.Usage)
	cmd.Annotations = map[string]string{
		group.Key: group.Workflows,
	}

	cmd.PersistentFlags().StringVarP(&profile, AWSProfileFlag, AWSProfileFlagShort, "", AWSProfileFlagDescription)

	return cmd
}
