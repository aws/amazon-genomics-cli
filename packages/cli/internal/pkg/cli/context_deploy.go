// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/slices"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	deployContextAllFlag        = "all"
	deployContextAllDescription = `Deploy all contexts in the project`
	deployContextDescription    = `Names of one or more contexts to deploy`
)

type ContextResult struct {
	Context string
	Err     error
}

type deployContextVars struct {
	contexts  []string
	deployAll bool
}

type deployContextOpts struct {
	deployContextVars
	ctxManager context.Interface
}

func newDeployContextOpts(vars deployContextVars) (*deployContextOpts, error) {
	return &deployContextOpts{
		deployContextVars: vars,
		ctxManager:        context.NewManager(profile),
	}, nil
}

func (o *deployContextOpts) Validate(contexts []string) error {
	o.contexts = append(o.contexts, contexts...)

	if (!o.deployAll && len(o.contexts) == 0) || (o.deployAll && len(o.contexts) > 0) {
		return fmt.Errorf("either an 'all' flag or a list of contexts must be provided, but not both")
	}

	if len(o.contexts) > 0 {
		if err := o.validateSuppliedContexts(o.contexts); err != nil {
			return err
		}
	} else {
		ctxList, err := o.ctxManager.List()
		if err != nil {
			return err
		}
		for contextName := range ctxList {
			o.contexts = append(o.contexts, contextName)
		}
	}

	return nil
}

func (o *deployContextOpts) validateSuppliedContexts(contextList []string) error {
	ctxList, err := o.ctxManager.List()
	if err != nil {
		return err
	}

	for _, contextName := range contextList {
		if _, ok := ctxList[contextName]; !ok {
			return fmt.Errorf("the provided context '%s' is not defined in the agc-project.yaml file", contextName)
		}
	}

	return nil
}

// Execute causes the specified context(s) to be deployed.
func (o *deployContextOpts) Execute() error {
	o.contexts = slices.DeDuplicateStrings(o.contexts)

	err := o.deployContexts()
	if err != nil {
		return err
	}

	log.Info().Msgf("Successfully deployed context(s) %s", o.contexts)
	return nil
}

func (o *deployContextOpts) deployContexts() error {
	progressResults := o.ctxManager.Deploy(o.contexts)
	var aggregateSuggestions []string

	var failedDeployments []context.ProgressResult
	for _, progressResult := range progressResults {
		if progressResult.Err != nil {
			failedDeployments = append(failedDeployments, progressResult)
		}
	}

	failedDeploymentsLength := len(failedDeployments)
	if failedDeploymentsLength > 0 {
		for i, failedDeployment := range failedDeployments {
			log.Error().Msgf("Failed to deploy context '%s'. Below is the log for that deployment", failedDeployment.Context)

			isLastDeployment := i == failedDeploymentsLength-1
			printErroredLogs(failedDeployment, isLastDeployment)

			var actionableError *actionableerror.Error
			if errors.As(failedDeployment.Err, &actionableError) {
				log.Error().Err(actionableError.Cause).Msgf(actionableError.Error())
				aggregateSuggestions = append(aggregateSuggestions, fmt.Sprintf("To resolve failure %d, try: %s", i+1, actionableError.SuggestedAction))
			} else {
				// An error occurred that we don't know how to deal with
				// already. We can't swallow it or it will be impossible for
				// the user to report or fix.
				aggregateSuggestions = append(aggregateSuggestions, fmt.Sprintf("To resolve failure %d, determine the cause of: %v", i+1, failedDeployment.Err))
			}
		}

		return actionableerror.New(fmt.Errorf("%d context deployment failures", failedDeploymentsLength), strings.Join(aggregateSuggestions, "\n"))
	}

	return nil
}

func printErroredLogs(failedDeployment context.ProgressResult, isLastDeployment bool) {
	outputsLength := len(failedDeployment.Outputs)
	if outputsLength == 0 {
		return
	}

	for j := 0; j < outputsLength-1; j++ {
		log.Error().Msg(failedDeployment.Outputs[j])
	}

	if isLastDeployment {
		log.Error().Msgf("%s \n\n\n", failedDeployment.Outputs[outputsLength-1])
	} else {
		log.Error().Msg(failedDeployment.Outputs[outputsLength-1])
	}
}

// BuildContextDeployCommand builds the command to deploy specified contexts in the current project.
func BuildContextDeployCommand() *cobra.Command {
	vars := deployContextVars{}
	cmd := &cobra.Command{
		Use:   "deploy {context_name ... | --all}",
		Short: "Deploy contexts in the current project",
		Long: `deploy is for deploying one or more contexts. 
It creates AGC resources in AWS.

` + DescribeOutput([]context.Detail{}),
		Example: `
/code agc context deploy context1 context2`,
		Args: cobra.ArbitraryArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newDeployContextOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(args); err != nil {
				return err
			}
			log.Info().Msgf("Deploying context(s)")
			err = opts.Execute()
			if err != nil {
				return clierror.New("context deploy", vars, err)
			}
			return nil
		}),
		ValidArgsFunction: NewContextAutoComplete().GetContextAutoComplete(),
	}
	cmd.Flags().BoolVar(&vars.deployAll, deployContextAllFlag, false, deployContextAllDescription)
	cmd.Flags().StringSliceVarP(&vars.contexts, contextFlag, contextFlagShort, nil, deployContextDescription)
	_ = cmd.RegisterFlagCompletionFunc(contextFlag, NewContextAutoComplete().GetContextAutoComplete())
	return cmd
}
