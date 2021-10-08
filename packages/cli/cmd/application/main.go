// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package main contains the root command.
package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/amazon-genomics-cli/cmd/application/template"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/aws/amazon-genomics-cli/internal/pkg/term/color"
	"github.com/aws/amazon-genomics-cli/internal/pkg/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const hugoFrontMatterTemplate = `---
date: %s
title: "%s"
slug: %s
---
`

type mainVars struct {
	docPath string
}

func init() {
	color.DisableColorBasedOnEnvVar()
	cobra.EnableCommandSorting = false // Maintain the order in which we add commands.
}

func main() {
	log.Logger = logging.ApplicationConsoleLogger()

	cmd := buildRootCmd()
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//BuildCommandDocsForHugo Generates markdown suitable for rendering by Hugo. Will only generate pages if 'dir' exists
func BuildCommandDocsForHugo(cmd *cobra.Command, dir string) error {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}

	frontMatterPrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(hugoFrontMatterTemplate, now, strings.Replace(base, "_", " ", -1), base)
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf("{{< relref \"%s\" >}}", strings.ToLower(base))
	}

	cmd.DisableAutoGenTag = true
	err := doc.GenMarkdownTreeCustom(cmd, dir, frontMatterPrepender, linkHandler)
	return err
}

func buildRootCmd() *cobra.Command {
	vars := mainVars{}
	cmd := &cobra.Command{
		Use:   "agc",
		Short: shortDescription,
		Example: `
  Displays the help menu for the specified sub-command.
  /code $ agc account --help`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// If we don't set a Run() function the help menu doesn't show up.
			// See https://github.com/spf13/cobra/issues/790
			if !logging.Verbose {
				// global level is trace by default so if not verbose we want info level
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if vars.docPath != "" {
				err := BuildCommandDocsForHugo(cmd, vars.docPath)
				if err != nil {
					return clierror.New("agc", args, err)
				}
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	// Sets version for --version flag. Version command gives more detailed
	// version information.
	cmd.Version = version.Version
	cmd.SetVersionTemplate("agc version: {{.Version}}\n")

	cmd.AddCommand(cli.BuildAccountCommand())
	cmd.AddCommand(cli.BuildProjectCommand())
	cmd.AddCommand(cli.BuildContextCommand())
	cmd.AddCommand(cli.BuildLogsCommand())
	cmd.AddCommand(cli.BuildWorkflowCommand())
	cmd.AddCommand(cli.BuildConfigureCommand())

	cmd.SetUsageTemplate(template.RootUsage)

	cmd.PersistentFlags().BoolVarP(&logging.Verbose, cli.VerboseFlag, cli.VerboseFlagShort, false, cli.VerboseFlagDescription)
	cmd.Flags().StringVar(&vars.docPath, "docs", "", "generate markdown documenting the CLI to the specified path")
	cmd.Flag("docs").Hidden = true

	return cmd
}
