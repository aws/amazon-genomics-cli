// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type contextStatusOpts struct {
	ctxManager context.Interface
}

func newContextStatusOpts() (*contextStatusOpts, error) {
	return &contextStatusOpts{
		ctxManager: context.NewManager(profile),
	}, nil
}

// Validate returns an error if the user's input is invalid.
func (o *contextStatusOpts) Validate() error {
	return nil
}

// Execute returns an array of status information strings about all context instances.
func (o *contextStatusOpts) Execute() ([]context.Instance, error) {
	return o.ctxManager.StatusList()
}

// BuildContextStatusCommand builds the command to show the status of a specific
// or for multiple context instances in the current project.
func BuildContextStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status for the deployed contexts in the project",
		Long: `status is for showing the status for the deployed contexts in the project. 

` + DescribeOutput([]context.Instance{}),
		Example: `
/code agc context status`,
		Args: cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newContextStatusOpts()
			if err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			contextInstances, err := opts.Execute()
			if err != nil {
				return clierror.New("context status", nil, err)
			}
			if len(contextInstances) > 0 {
				format.Default.Write(contextInstances)
			} else {
				log.Info().Msg("There are no contexts deployed.")
			}
			return nil
		}),
	}
	return cmd
}
