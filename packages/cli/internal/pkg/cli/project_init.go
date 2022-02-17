// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"sort"

	"github.com/aws/amazon-genomics-cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/group"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
	"github.com/spf13/cobra"
)

const initContextName = "ctx1"

const (
	projectInitWorkflowTypeName      = "workflow-type"
	projectInitWorkflowTypeNameShort = "w"
)

var (
	workflowTypeToEngineMap = map[string]string{
		"nextflow":  "nextflow",
		"wdl":       "cromwell",
		"snakemake": "snakemake",
	}
	supportedWorkflowTypes []string
)

func init() {
	var workflowTypes []string
	for name := range workflowTypeToEngineMap {
		workflowTypes = append(workflowTypes, name)
	}
	sort.Strings(workflowTypes)
	supportedWorkflowTypes = workflowTypes
}

type initProjectVars struct {
	ProjectName  string
	workflowType string
}

type initProjectOpts struct {
	initProjectVars
	projectClient storage.ProjectClient
}

func (o *initProjectOpts) validateWorkflowType() error {
	if _, ok := workflowTypeToEngineMap[o.workflowType]; !ok {
		return fmt.Errorf("invalid workflow type supplied: '%s'. Supported workflow types are: %v", o.workflowType, supportedWorkflowTypes)
	}
	return nil
}

func (o *initProjectOpts) generateEngine() []spec.Engine {
	return []spec.Engine{{Type: o.workflowType, Engine: getEngineForWorkflowType(o.workflowType)}}
}

func getEngineForWorkflowType(workflowType string) string {
	return workflowTypeToEngineMap[workflowType]
}

func getProjectInitWorkflowTypeNameDescription() string {
	return fmt.Sprintf("uses the specified workflow type for the default context. Valid values include %v", supportedWorkflowTypes)
}

func newInitProjectOpts(vars initProjectVars) (*initProjectOpts, error) {
	projectClient, err := storage.NewProjectClientInCurrentDir()
	if err != nil {
		return nil, err
	}

	return &initProjectOpts{
		initProjectVars: vars,
		projectClient:   projectClient,
	}, nil
}

// Validate returns an error if specified project name is invalid
func (o *initProjectOpts) Validate() error {
	if err := validateProjectName(o.ProjectName); err != nil {
		return err
	}
	return o.validateProject()
}

// Execute creates a new empty project specification.
func (o *initProjectOpts) Execute() error {
	newProject := o.createInitialProject()
	return o.projectClient.Write(newProject)
}

func (o *initProjectOpts) createInitialProject() spec.Project {
	return spec.Project{
		Name:          o.ProjectName,
		SchemaVersion: spec.LatestVersion,
		Contexts: map[string]spec.Context{
			initContextName: {RequestSpotInstances: false, Engines: o.generateEngine()},
		},
	}
}

func (o *initProjectOpts) validateProject() error {
	if o.workflowType == "" {
		return fmt.Errorf("please specify a workflow type with the --%s flag", projectInitWorkflowTypeName)
	}
	if err := o.validateWorkflowType(); err != nil {
		return err
	}
	isInitialized, err := o.projectClient.IsInitialized()
	if err != nil {
		return err
	}
	if isInitialized {
		return fmt.Errorf("project specification already exists in folder '%s'", o.projectClient.GetLocation())
	}
	return nil
}

func BuildProjectInitCommand() *cobra.Command {
	vars := initProjectVars{}
	cmd := &cobra.Command{
		Use:   "init project_name --workflow-type {wdl|nextflow|snakemake}",
		Short: "Initialize current directory with a new empty AGC project for a particular workflow type.",
		Long: `Initialize current directory with a new empty AGC project for a particular workflow type.
Project specification file 'agc-project.yaml' will be created in the current directory.`,
		Example: `
Initialize a new project named "myProject".
/code $ agc project init myProject --workflow-type my_workflow_type`,
		Args: cobra.ExactArgs(1),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			vars.ProjectName = args[0]
			opts, err := newInitProjectOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			if err := opts.Execute(); err != nil {
				return clierror.New("project init", vars, err)
			}
			return nil
		}),
		Annotations: map[string]string{
			group.Key: group.GettingStarted,
		},
	}
	cmd.SetUsageTemplate(template.Usage)

	cmd.Flags().StringVarP(&vars.workflowType, projectInitWorkflowTypeName, projectInitWorkflowTypeNameShort, "", getProjectInitWorkflowTypeNameDescription())
	return cmd
}
