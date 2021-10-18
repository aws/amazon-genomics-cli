// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"sort"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type listContextVars struct{}

type listContextOpts struct {
	listContextVars
	ctxManager context.Interface
}

func newListContextOpts(vars listContextVars) (*listContextOpts, error) {
	return &listContextOpts{
		listContextVars: vars,
		ctxManager:      context.NewManager(profile),
	}, nil
}

func (o *listContextOpts) Validate() error {
	return nil
}

// Execute returns a collection of contexts defined in project specification
func (o *listContextOpts) Execute() ([]types.ContextSummary, error) {
	contexts, err := o.ctxManager.List()
	if err != nil {
		return nil, err
	}

	var contextNames []types.ContextSummary
	for name := range contexts {
		if len(contexts[name].Engines) != 1 {
			return nil, actionableerror.New(fmt.Errorf("context '%s' does not have a valid engine declaration", name), "please validate your project yaml with 'agc project validate'")
		}
		contextNames = append(contextNames, types.ContextSummary{
			Name:       name,
			EngineName: contexts[name].Engines[0].Engine,
		})
	}

	sortContextNames(contextNames)
	return contextNames, nil
}

func sortContextNames(contextNames []types.ContextSummary) {
	sort.Slice(contextNames, func(i, j int) bool {
		return contextNames[i].Name < contextNames[j].Name
	})
}

func BuildContextListCommand() *cobra.Command {
	vars := listContextVars{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List contexts in the project",
		Long: `list is for showing a combined list of contexts specified in
the project specification.

` + DescribeOutput([]types.ContextSummary{}),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newListContextOpts(vars)
			if err != nil {
				return err
			}
			log.Info().Msgf("Listing contexts.")
			if err := opts.Validate(); err != nil {
				return err
			}
			contexts, err := opts.Execute()
			if err != nil {
				return clierror.New("context list", vars, err)
			}
			format.Default.Write(contexts)
			return nil
		}),
	}
	return cmd
}
