// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/actionable"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	appTagKey                      = "application-name"
	appTagValue                    = "agc"
	deactivateForceFlag            = "force"
	deactivateForceShortFlag       = "f"
	deactivateForceFlagDescription = `Force account deactivation by removing all resources associated with AGC.
This includes project and context resources, even if they are running workflows.
If not specified, only the core resources will be deleted if possible.`
)

type accountDeactivateVars struct {
	force bool
}

type accountDeactivateOpts struct {
	accountDeactivateVars
	stacks    []cfn.Stack
	cfnClient cfn.Interface
}

func newAccountDeactivateOpts(vars accountDeactivateVars) (*accountDeactivateOpts, error) {
	return &accountDeactivateOpts{
		accountDeactivateVars: vars,
		cfnClient:             aws.CfnClient(profile),
	}, nil
}
func (o *accountDeactivateOpts) LoadStacks() error {
	stacks, err := o.getApplicationStacks()
	if err != nil {
		return err
	}
	o.stacks = stacks
	return nil
}

func (o *accountDeactivateOpts) Validate() error {
	if !o.force && len(o.stacks) > 1 {
		return actionable.NewError(
			errors.New("one or more contexts are still deployed"),
			"use --force to destroy deployed contexts as well",
		)
	}
	return nil
}

func (o *accountDeactivateOpts) Execute() error {
	stackDeletionTrackers := make(map[string]chan cfn.DeletionResult)
	for _, stack := range o.stacks {
		log.Debug().Msgf("Deleting stack '%s'", stack.Name)
		tracker, err := o.cfnClient.DeleteStack(stack.Id)
		if err != nil {
			return err
		}
		stackDeletionTrackers[stack.Name] = tracker
	}

	for stackName, tracker := range stackDeletionTrackers {
		deletionResult := <-tracker
		if deletionResult.Error != nil {
			return fmt.Errorf("failed to delete stack '%s: %w", stackName, deletionResult.Error)
		}
		log.Debug().Msgf("Stack '%s' has been successfully deleted!", stackName)
	}

	return nil
}

func (o *accountDeactivateOpts) getApplicationStacks() ([]cfn.Stack, error) {
	stacks, err := o.cfnClient.ListStacks(regexp.MustCompile(`^Agc-.*$`), cfn.ActiveStacksFilter)
	if err != nil {
		return nil, err
	}

	var filteredStacks []cfn.Stack
	for _, stack := range stacks {
		tags, err := o.cfnClient.GetStackTags(stack.Id)
		if err != nil {
			return nil, err
		}
		if tags[appTagKey] == appTagValue {
			filteredStacks = append(filteredStacks, stack)
		}
	}

	return filteredStacks, nil
}

// BuildAccountDeactivateCommand builds the command for deactivating AGC in an AWS account.
func BuildAccountDeactivateCommand() *cobra.Command {
	vars := accountDeactivateVars{}
	cmd := &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate AGC in an AWS account.",
		Long: `Deactivate AGC in an AWS account.
AGC will use your default AWS credentials to remove all core AWS resources
it has created in that account and region. Deactivation may take up to 5 minutes to complete and return.
Buckets and logs will be preserved.`,
		Example: `
Deactivate AGC in your AWS account.
/code $ agc account deactivate`,
		Args: cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newAccountDeactivateOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.LoadStacks(); err != nil {
				return clierror.New("account deactivate", vars, err)
			}
			if err := opts.Validate(); err != nil {
				return clierror.New("account deactivate", vars, err)
			}
			log.Info().Msgf("Deactivating AGC. Deactivation may take up to 5 minutes to complete and return.")
			if err := opts.Execute(); err != nil {
				return clierror.New("account deactivate", vars, err)
			}
			return nil
		}),
	}
	cmd.Flags().BoolVarP(&vars.force, deactivateForceFlag, deactivateForceShortFlag, false, deactivateForceFlagDescription)
	return cmd
}
