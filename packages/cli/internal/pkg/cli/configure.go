// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/group"
	"github.com/spf13/cobra"
)

// BuildConfigureCommand builds the top level configure command
func BuildConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "configure",
		Short: `Commands for configuration.
Configuration is stored per user.`,
		Long: `Commands for configuration.
Configure local settings and preferences to customize the CLI experience.`,
	}

	cmd.AddCommand(BuildConfigureEmailCommand())
	cmd.AddCommand(BuildDescribeShowCommand())

	cmd.SetUsageTemplate(template.Usage)
	cmd.Annotations = map[string]string{
		"group": group.Settings,
	}

	cmd.PersistentFlags().StringVarP(&profile, AWSProfileFlag, AWSProfileFlagShort, "", AWSProfileFlagDescription)

	return cmd
}
