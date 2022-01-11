// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type describeProjectVars struct {
}

type describeProjectOpts struct {
	projectClient storage.ProjectClient
	describeProjectVars
}

func newDescribeProjectOpts(vars describeProjectVars) (*describeProjectOpts, error) {
	projectClient, err := storage.NewProjectClient()
	if err != nil {
		return nil, err
	}
	return &describeProjectOpts{
		describeProjectVars: vars,
		projectClient:       projectClient,
	}, nil
}

// Validate returns an error if specified project is not initialized
func (o *describeProjectOpts) Validate() error {
	isInitialized, err := o.projectClient.IsInitialized()
	if err != nil {
		return err
	}
	if !isInitialized {
		return fmt.Errorf("project at location '%s' is not initialized", o.projectClient.GetLocation())
	}
	return nil
}

func (o *describeProjectOpts) Execute() (types.Project, error) {
	projectSpec, err := o.projectClient.Read()
	if err != nil {
		return types.Project{}, err
	}
	dataRefs := buildDataRefs(projectSpec)
	if err != nil {
		return types.Project{}, err
	}
	return types.Project{
		Name: projectSpec.Name,
		Data: dataRefs,
	}, nil
}

func buildDataRefs(projectSpec spec.Project) []types.Data {
	var dataList []types.Data
	for _, data := range projectSpec.Data {
		dataType := types.Data{Location: data.Location, ReadOnly: data.ReadOnly}
		dataList = append(dataList, dataType)
	}
	return dataList
}

func buildProjectDescribeCommand() *cobra.Command {
	vars := describeProjectVars{}
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a project",
		Long: `Describe is for describing the current project specification. The current project specification is determined to be the agc-project.yaml file in the current working directory or a parent of the current directory

` + DescribeOutput(types.Project{}),
		Example: `
/code agc project describe`,
		Args: cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newDescribeProjectOpts(vars)
			if err != nil {
				return err
			}
			log.Info().Msgf("Describing current project")
			if err := opts.Validate(); err != nil {
				return err
			}
			project, err := opts.Execute()
			if err != nil {
				return clierror.New("project describe", vars, err)
			}
			format.Default.Write(project)
			return nil
		}),
	}
	return cmd
}
