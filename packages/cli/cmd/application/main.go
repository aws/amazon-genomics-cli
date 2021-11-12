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
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/config"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/aws/amazon-genomics-cli/internal/pkg/storage"
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
type formatVars struct {
	format string
}

var newConfigClient = func() (storage.ConfigClient, error) {
	return config.NewConfigClient()
}

const (
	defaultFormat = "text"
)

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
	formatVars := formatVars{}
	cmd := &cobra.Command{
		Use:   "agc",
		Short: shortDescription,
		Example: `
  Displays the help menu for the specified sub-command.
  /code $ agc account --help`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLoggingLevel()
			setFormatter(formatVars)
			checkCliVersion()
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
	cmd.PersistentFlags().StringVar(&formatVars.format, cli.FormatFlag, "", cli.FormatFlagDescription)
	cmd.Flags().StringVar(&vars.docPath, "docs", "", "generate markdown documenting the CLI to the specified path")
	cmd.Flag("docs").Hidden = true

	return cmd
}

func ValidateFormat(f format.FormatterType) error {
	if err := f.ValidateFormatter(); err != nil {
		return err
	}
	return nil
}

func setFormatter(f formatVars) string {
	configClient, err := newConfigClient()
	if err != nil {
		log.Error().Err(err)
		return ""
	}
	if f.format == "" {
		f.format = defaultFormat
		configFormat, err := configClient.GetFormat()
		if err != nil {
			log.Error().Err(err)
		} else {
			f.format = configFormat
		}
	}
	if err := ValidateFormat(format.FormatterType(f.format)); err != nil {
		fmt.Println(err.Error())
	}
	format.SetFormatter(format.FormatterType(f.format))
	return f.format
}

func checkCliVersion() {
	log.Debug().Msg("Checking AGC version...")
	result, err := version.Check()
	if err != nil {
		log.Debug().Msgf("version check failed: %v", err)
		return
	}
	if result.LatestVersion != result.CurrentVersion {
		log.Info().Msgf("New version of agc available. Current version is '%s'. The latest version is '%s'", result.CurrentVersion, result.LatestVersion)
	}
	if result.CurrentVersionDeprecated {
		log.Warn().Msgf("The current version was deprecated with message: %s. Please consider upgrading Amazon Genomics CLI.", strings.TrimSpace(result.CurrentVersionDeprecationMessage))
	}
	if len(result.NewerVersionHighlights) > 0 {
		log.Info().Msgf("Highlights from newer versions:")
		for _, msg := range result.NewerVersionHighlights {
			log.Info().Msgf("\t%s", strings.TrimSpace(msg))
		}
	}
}

func setLoggingLevel() {
	if !logging.Verbose {
		// global level is trace by default so if not verbose we want info level
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
