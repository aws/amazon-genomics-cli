// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

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
		ctxList, err := o.ctxManagerFactory().List()
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
	ctxList, err := o.ctxManagerFactory().List()
	if err != nil {
		return err
	}

	for _, context := range contextList {
		if _, ok := ctxList[context]; !ok {
			return fmt.Errorf("the provided context '%s' is not defined in the agc-project.yaml file", context)
		}
	}

	return nil
}

// Execute causes the specified context(s) to be deployed.
func (o *deployContextOpts) Execute() ([]context.Detail, error) {
	o.contexts = slices.DeDuplicateStrings(o.contexts)

	err := o.deployContexts()
	if err != nil {
		return nil, err
	}

	contextDetails, err := o.validateDeploymentResults()
	if err != nil {
		return nil, err
	}

	sortContextDetails(contextDetails)
	return contextDetails, nil
}

func (o *deployContextOpts) deployContexts() error {
	progressResults := o.ctxManagerFactory().Deploy(o.contexts)
	aggregateSuggestions := make([]string, 0)

	failedDeployments := make([]context.ProgressResult, 0)
	for _, progressResult := range progressResults {
		if progressResult.Err != nil {
			failedDeployments = append(failedDeployments, progressResult)
		}
	}

	failedDeploymentsLength := len(failedDeployments)
	if failedDeploymentsLength > 0 {
		for i, failedDeployment := range failedDeployments {
			log.Error().Msgf("Failed to deploy context '%s'. Below is the log for that deployment", failedDeployment.Context)

			outputsLength := len(failedDeployment.Outputs)
			for i := 0; i < outputsLength-1; i++ {
				log.Error().Msg(failedDeployment.Outputs[i])
			}

			if i < failedDeploymentsLength-1 {
				log.Error().Msgf("%s \n\n\n", failedDeployment.Outputs[outputsLength-1])
			} else {
				log.Error().Msg(failedDeployment.Outputs[outputsLength-1])
			}

			var actionableError *actionableerror.Error
			ok := errors.As(failedDeployment.Err, &actionableError)
			if ok {
				log.Error().Err(actionableError.Cause).Msgf(actionableError.Error())
				aggregateSuggestions = append(aggregateSuggestions, actionableError.SuggestedAction)
			}
		}

		return actionableerror.New(fmt.Errorf("one or more contexts failed to deploy"), strings.Join(aggregateSuggestions, ", "))
	}

	return nil
}

func (o *deployContextOpts) validateDeploymentResults() ([]context.Detail, error) {
	contextDetails, deploymentHasErrors := make([]context.Detail, len(o.contexts)), false
	aggregateSuggestions := make([]string, 0)
	var wait sync.WaitGroup
	wait.Add(len(o.contexts))

	for i, contextName := range o.contexts {
		go func(ctxManager context.Interface, i int, contextName string) {
			info, err := ctxManager.Info(contextName)
			contextDetails[i] = info
			if err != nil {
				var actionableError *actionableerror.Error
				ok := errors.As(err, &actionableError)
				if ok {
					log.Error().Err(actionableError.Cause).Msgf(actionableError.Error())
					aggregateSuggestions = append(aggregateSuggestions, actionableError.SuggestedAction)
				} else {
					log.Error().Err(err).Msgf("failed to deploy context '%s'", contextName)
				}
				deploymentHasErrors = true
			}
			wait.Done()
		}(o.ctxManagerFactory(), i, contextName)
	}

	wait.Wait()
	if deploymentHasErrors {
		aggregateSuggestions = slices.DeDuplicateStrings(aggregateSuggestions)
		return nil, actionableerror.New(fmt.Errorf("one or more contexts failed to deploy"), strings.Join(aggregateSuggestions, ", "))
	}

	return contextDetails, nil
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
