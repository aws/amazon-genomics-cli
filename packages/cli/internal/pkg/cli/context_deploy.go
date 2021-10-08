// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/slices"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	deployContextAllFlag        = "all"
	deployContextAllDescription = `Deploy all contexts in the project`
	deployContextDescription    = `Names of one or more contexts to deploy`
)

type deployResult struct {
	contextName string
	info        context.Detail
	err         error
}

type deployContextVars struct {
	contexts  []string
	deployAll bool
}

type deployContextOpts struct {
	deployContextVars
	ctxManagerFactory func() context.Interface
}

func newDeployContextOpts(vars deployContextVars) (*deployContextOpts, error) {
	return &deployContextOpts{
		deployContextVars: vars,
		ctxManagerFactory: func() context.Interface { return context.NewManager(profile) },
	}, nil
}

func (o *deployContextOpts) Validate() error {
	if (!o.deployAll && len(o.contexts) == 0) || (o.deployAll && len(o.contexts) > 0) {
		return fmt.Errorf("one of either the 'context' or 'all' flag is required")
	}
	return nil
}

// Execute causes the specified context(s) to be deployed.
func (o *deployContextOpts) Execute() ([]context.Detail, error) {
	if o.deployAll {
		ctxList, err := o.ctxManagerFactory().List()
		if err != nil {
			return nil, err
		}
		for contextName := range ctxList {
			o.contexts = append(o.contexts, contextName)
		}
	}

	o.contexts = slices.DeDuplicateStrings(o.contexts)
	results := o.deployContexts(o.contexts)
	contextDetails := make([]context.Detail, len(results))
	hasErrors := false
	aggregateSuggestions := make([]string, 0)
	for i, result := range results {
		if result.err != nil {
			var actionableError *actionableerror.Error
			ok := errors.As(result.err, &actionableError)
			if ok {
				log.Error().Err(actionableError.Cause).Msgf(actionableError.Error())
				aggregateSuggestions = append(aggregateSuggestions, actionableError.SuggestedAction)
			} else {
				log.Error().Err(result.err).Msgf("failed to deploy context '%s'", result.contextName)
			}
			hasErrors = true
		}
		contextDetails[i] = result.info
	}
	if hasErrors {
		aggregateSuggestions = slices.DeDuplicateStrings(aggregateSuggestions)
		return nil, actionableerror.New(fmt.Errorf("one or more contexts failed to deploy"), strings.Join(aggregateSuggestions, ", "))
	}
	sortContextDetails(contextDetails)
	return contextDetails, nil
}

func (o *deployContextOpts) deployContexts(contexts []string) []deployResult {
	results := make([]deployResult, len(contexts))
	for i, contextName := range contexts {
		log.Debug().Msgf("Deploying context '%s'", contextName)
		// TODO: Run in parallel once CDK resolves race condition causing context bleed
		//       https://github.com/aws/aws-cdk/issues/14350
		func(ctxManager context.Interface, i int, contextName string) {
			_ = ctxManager.Deploy(contextName, true)
			info, err := ctxManager.Info(contextName)
			results[i] = deployResult{contextName: contextName, info: info, err: err}
		}(o.ctxManagerFactory(), i, contextName)
	}
	return results
}

func sortContextDetails(contextDetails []context.Detail) {
	sort.Slice(contextDetails, func(i, j int) bool {
		return contextDetails[i].Name < contextDetails[j].Name
	})
}

// BuildContextDeployCommand builds the command to deploy specified contexts in the current project.
func BuildContextDeployCommand() *cobra.Command {
	vars := deployContextVars{}
	cmd := &cobra.Command{
		Use:   "deploy {-c context_name ... | --all}",
		Short: "Deploy contexts in the current project",
		Long: `deploy is for deploying one or more contexts. 
It creates AGC resources in AWS.

` + DescribeOutput([]context.Detail{}),
		Example: `
/code agc context deploy -c context1 -c context2`,
		Args: cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newDeployContextOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			log.Info().Msgf("Deploying context(s)")
			contextInfo, err := opts.Execute()
			if err != nil {
				return clierror.New("context deploy", vars, err)
			}
			format.Default.Write(contextInfo)
			return nil
		}),
	}
	cmd.Flags().BoolVar(&vars.deployAll, deployContextAllFlag, false, deployContextAllDescription)
	cmd.Flags().StringSliceVarP(&vars.contexts, contextFlag, contextFlagShort, nil, deployContextDescription)
	_ = cmd.RegisterFlagCompletionFunc(contextFlag, NewContextAutoComplete().GetContextAutoComplete())
	return cmd
}
