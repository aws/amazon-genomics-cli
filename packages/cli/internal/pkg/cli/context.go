// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/group"
	"github.com/spf13/cobra"
)

// BuildContextCommand builds the top level context command and related subcommands.
func BuildContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "context",
		Short: `Commands for contexts.
Contexts specify workflow engines and computational fleets to use when running a workflow.`,
		Long: `Commands for contexts.
Contexts specify workflow engines and computational fleets to use when running a workflow.
Users can quickly switch between infrastructure configurations by specifying a 
particular context.`,
	}

	cmd.AddCommand(BuildContextDescribeCommand())
	cmd.AddCommand(BuildContextListCommand())
	cmd.AddCommand(BuildContextDeployCommand())
	cmd.AddCommand(BuildContextDestroyCommand())
	cmd.AddCommand(BuildContextStatusCommand())

	cmd.SetUsageTemplate(template.Usage)
	cmd.Annotations = map[string]string{
		group.Key: group.Contexts,
	}

	cmd.PersistentFlags().StringVarP(&profile, AWSProfileFlag, AWSProfileFlagShort, "", AWSProfileFlagDescription)

	return cmd
}
