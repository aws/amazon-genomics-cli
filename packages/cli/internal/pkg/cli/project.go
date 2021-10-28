// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"regexp"

	"github.com/aws/amazon-genomics-cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/group"
	"github.com/spf13/cobra"
)

// ProjectInfo defines the metadata that is kept for a project.
type ProjectInfo struct {
	ProjectName    string
	Workspace      string
	OutputLocation string
}

func validateProjectName(projectName string) error {
	if projectName == "" {
		return fmt.Errorf("missing project name")
	}
	isAlphaNumeric := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNumeric(projectName) {
		return fmt.Errorf("%s has non-alpha-numeric characters in it", projectName)
	}
	return nil
}

// BuildProjectCommand builds the top level project command and related subcommands.
func BuildProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: `Commands to interact with projects.`,
		Long: `Commands to interact with projects.
A project is a local configuration file that describes the workflows, data,
and contexts you are working with.`,
	}

	cmd.AddCommand(BuildProjectInitCommand())
	cmd.AddCommand(buildProjectDescribeCommand())
	cmd.AddCommand(buildProjectValidateCommand())

	cmd.SetUsageTemplate(template.Usage)
	cmd.Annotations = map[string]string{
		group.Key: group.Projects,
	}

	return cmd
}
