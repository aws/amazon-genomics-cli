// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/group"
	"github.com/spf13/cobra"
)

// BuildLogsCommand builds the top level logs command and related subcommands.
func BuildLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: `Commands for various logs.`,
		Long: `Commands for various logs.
Logs can currently be listed for workflows, workflow engines, 
and various AGC infrastructure parts.
You can also show the content of any CloudWatch log stream that you have
access rights to.`,
	}

	cmd.AddCommand(BuildLogsWorkflowCommand())
	cmd.AddCommand(BuildLogsEngineCommand())
	cmd.AddCommand(BuildLogsAdapterCommand())
	cmd.AddCommand(BuildLogsAccessCommand())

	cmd.SetUsageTemplate(template.Usage)
	cmd.Annotations = map[string]string{
		group.Key: group.Logs,
	}

	cmd.PersistentFlags().StringVarP(&profile, AWSProfileFlag, AWSProfileFlagShort, "", AWSProfileFlagDescription)

	return cmd
}
