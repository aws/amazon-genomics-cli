// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type validateProjectOpts struct {
	projectClient storage.ProjectClient
}

func newValidateProjectOpts() (*validateProjectOpts, error) {
	projectClient, err := storage.NewProjectClient()
	if err != nil {
		return nil, err
	}
	return &validateProjectOpts{
		projectClient: projectClient,
	}, nil
}

func (o *validateProjectOpts) Execute() error {
	_, err := o.projectClient.Read()
	return err
}

func buildProjectValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate an agc-project.yaml file",
		Long: `Determines if the current project specification follows the required format and lists any syntax errors. 
The current project specification is determined to be the agc-project.yaml file in the current working directory or a parent of the current directory
` + DescribeOutput(types.Project{}),
		Example: `
/code agc project describe`,
		Args: cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newValidateProjectOpts()
			if err != nil {
				return err
			}
			log.Info().Msgf("Validating specification at project root: '%s'", opts.projectClient.GetLocation())
			err = opts.Execute()
			if err != nil {
				return clierror.New("project describe", "", err)
			}
			log.Info().Msgf("No errors found.")
			return nil
		}),
	}
	return cmd
}
