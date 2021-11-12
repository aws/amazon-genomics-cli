// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"net/mail"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const configureEmailCommand = "configure email"

type emailContextVars struct {
	userEmailAddress string
}
type emailContextOpts struct {
	emailContextVars
	configClient config.ConfigClient
}

func newEmailContextOpts(vars emailContextVars) (*emailContextOpts, error) {
	return &emailContextOpts{
		emailContextVars: vars,
	}, nil
}

func (o *emailContextOpts) Validate() error {
	_, err := mail.ParseAddress(o.userEmailAddress)
	return err
}

// Execute returns a context definition for the specified name.
func (o *emailContextOpts) Execute() error {
	err := o.configClient.SetUserEmailAddress(o.emailContextVars.userEmailAddress)
	if err != nil {
		return err
	}

	return nil
}
func BuildConfigureEmailCommand() *cobra.Command {
	vars := emailContextVars{}
	cmd := &cobra.Command{
		Use:   "email user_email_address",
		Short: "Sets user email address to be used to tag AGC resources created in the account",
		Args:  cobra.ExactArgs(1),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				vars.userEmailAddress = args[0]
				opts, err := newEmailContextOpts(vars)
				if err != nil {
					return err
				}
				configClient, err := config.NewConfigClient()
				if err != nil {
					return clierror.New(configureEmailCommand, vars, err)
				}
				opts.configClient = configClient
				log.Info().Msgf("Setting user email address to: '%s'", opts.userEmailAddress)
				if err := opts.Validate(); err != nil {
					return err
				}
				err = opts.Execute()
				if err != nil {
					return clierror.New(configureEmailCommand, vars, err)
				}
			}

			return nil
		}),
	}
	return cmd
}
