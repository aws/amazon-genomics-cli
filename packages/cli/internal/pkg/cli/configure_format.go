// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const configureFormatCommand = "configure format"

type formatContextVars struct {
	format string
}
type formatContextOpts struct {
	formatContextVars
	configClient storage.ConfigClient
}

func newFormatContextOpts(vars formatContextVars) (*formatContextOpts, error) {
	return &formatContextOpts{
		formatContextVars: vars,
	}, nil
}

func (o *formatContextOpts) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("a single format value must be provided")
	}
	format := format.FormatterType(o.formatContextVars.format)
	if err := format.ValidateFormatter(); err != nil {
		return err
	}
	return nil
}

func (o *formatContextOpts) Execute() error {
	err := o.configClient.SetFormat(o.formatContextVars.format)
	if err != nil {
		return err
	}

	return nil
}
func BuildConfigureFormatCommand() *cobra.Command {
	vars := formatContextVars{}
	cmd := &cobra.Command{
		Use:   "format output_format",
		Short: "Sets default format option for output display of AGC commands. Valid format options are 'text' and 'table'",
		Args:  cobra.ArbitraryArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newFormatContextOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(args); err != nil {
				return err
			}
			vars.format = args[0]
			configClient, err := config.NewConfigClient()
			if err != nil {
				return clierror.New(configureFormatCommand, vars, err)
			}
			opts.configClient = configClient
			log.Info().Msgf("Setting default format to: '%s'", opts.format)
			err = opts.Execute()
			if err != nil {
				return clierror.New(configureFormatCommand, vars, err)
			}
			return nil
		}),
	}
	return cmd
}
