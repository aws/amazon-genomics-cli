// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type showContextOpts struct {
	configClient storage.ConfigClient
}

func newConfigureDescribeContextOpts() (*showContextOpts, error) {
	configClient, err := config.NewConfigClient()
	if err != nil {
		return nil, err
	}

	return &showContextOpts{configClient}, nil
}

func (o *showContextOpts) Validate() error {
	return nil
}

// Execute returns current user specific configuration snapshot
func (o *showContextOpts) Execute() (config.Config, error) {
	return o.configClient.Read()
}

func BuildDescribeShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows current configuration of the AGC setup for current user",
		Long:  "Running this command reads current configuration file for AGC and prints out it content\n" + DescribeOutput(config.Config{}),
		Args:  cobra.ExactArgs(0),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newConfigureDescribeContextOpts()
			if err != nil {
				return err
			}

			log.Info().Msgf("Reading user specific configuration")
			if err := opts.Validate(); err != nil {
				return err
			}
			configuration, err := opts.Execute()
			if err != nil {
				return clierror.New("configure describe", nil, err, "you may need to configure your account")
			}
			format.Default.Write(configuration)

			return nil
		}),
	}
	return cmd
}
