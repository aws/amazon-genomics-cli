// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/group"
	"github.com/spf13/cobra"
)

// BuildAccountCommand builds the top level account command and related subcommands.
func BuildAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "account",
		Short: `Commands for AWS account setup.
Install or remove AGC from your account.`,
		Long: `Commands for AWS account setup.
AGC requires core infrastructure to be running in an account to function.
These commands should be used to install or remove AGC from your AWS account.`,
	}

	cmd.AddCommand(BuildAccountActivateCommand())
	cmd.AddCommand(BuildAccountDeactivateCommand())

	cmd.SetUsageTemplate(template.Usage)
	cmd.Annotations = map[string]string{
		"group": group.GettingStarted,
	}

	cmd.PersistentFlags().StringVarP(&profile, AWSProfileFlag, AWSProfileFlagShort, "", AWSProfileFlagDescription)

	return cmd
}
