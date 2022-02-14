// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	destroyContextAllFlag          = "all"
	destroyContextAllDescription   = `Destroy all contexts in the project`
	destroyContextForceFlag        = "force"
	destroyContextForceDescription = "Destroy context and stop running workflows within context"
	destroyContextDescription      = `Names of one or more contexts to destroy`
)

type destroyResult struct {
	contextName string
	err         error
}

type destroyContextVars struct {
	contexts     []string
	destroyAll   bool
	destroyForce bool
}

type destroyContextOpts struct {
	destroyContextVars
	ctxManagerFactory func() context.Interface
	wfsManager        func() workflow.Interface
}

func newDestroyContextOpts(vars destroyContextVars) (*destroyContextOpts, error) {
	contextOpts := &destroyContextOpts{
		destroyContextVars: vars,
		ctxManagerFactory:  func() context.Interface { return context.NewManager(profile) },
		wfsManager:         func() workflow.Interface { return workflow.NewManager(profile) },
	}

	return contextOpts, nil
}

func (o *destroyContextOpts) Validate(contexts []string) error {
	o.contexts = append(o.contexts, contexts...)

	if (!o.destroyAll && len(o.contexts) == 0) || (o.destroyAll && len(o.contexts) > 0) {
		return fmt.Errorf("either an 'all' flag or a list of contexts must be provided, but not both")
	}

	err := o.validateContexts()
	if err != nil {
		return err
	}

	wfsManager := o.wfsManager()
	for _, ctx := range o.contexts {
		workflows, err := wfsManager.StatusWorkflowByContext(ctx, workflowMaxAllowedInstance)
		if err != nil {
			return err
		}
		for _, wf := range workflows {
			if wf.IsInstanceRunning() {
				if !o.destroyForce {
					return fmt.Errorf("context '%s' contains running workflows. "+
						"Please stop all workflows before destroying context", ctx)
				} else {
					wfsManager.StopWorkflowInstance(wf.Id)
				}
			}
		}
	}
	return nil
}

// Execute causes the specified context(s) to be destroyed.
func (o *destroyContextOpts) Execute() error {
	results := o.destroyContexts(o.contexts)
	hasErrors := false
	for _, result := range results {
		if result.err != nil {
			log.Error().Err(result.err).Msgf("failed to destroy context '%s'", result.contextName)
			hasErrors = true
		}
	}
	if hasErrors {
		return fmt.Errorf("one or more contexts failed to be destroyed")
	}

	return nil
}

func (o *destroyContextOpts) validateContexts() error {
	ctxList, err := o.ctxManagerFactory().List()
	if err != nil {
		return err
	}
	if o.destroyAll {
		for contextName := range ctxList {
			o.contexts = append(o.contexts, contextName)
		}
	}

	for _, context := range o.contexts {
		if _, ok := ctxList[context]; !ok {
			return fmt.Errorf("the provided context '%s' is not defined in the agc-project.yaml file", context)
		}
	}

	return nil
}

func (o *destroyContextOpts) destroyContexts(contexts []string) []destroyResult {
	results := make([]destroyResult, len(contexts))
	destroyInfos := o.ctxManagerFactory().Destroy(contexts)

	for index, destroyInfo := range destroyInfos {
		results[index] = destroyResult{
			err:         destroyInfo.Err,
			contextName: destroyInfo.Context,
		}
	}
	return results
}

func BuildContextDestroyCommand() *cobra.Command {
	vars := destroyContextVars{}
	cmd := &cobra.Command{
		Use:   "destroy {context_name ... | --all}",
		Short: "Destroy contexts in the current project.",
		Long: `destroy is for destroying one or more contexts. 
It destroys AGC resources in AWS.`,
		Example: `
/code agc context destroy context1 context2`,
		Args: cobra.ArbitraryArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newDestroyContextOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(args); err != nil {
				return err
			}
			log.Info().Msgf("Destroying context(s)'")
			err = opts.Execute()
			if err != nil {
				return clierror.New("context destroy", vars, err)
			}
			return nil
		}),
	}
	cmd.Flags().BoolVar(&vars.destroyAll, destroyContextAllFlag, false, destroyContextAllDescription)
	cmd.Flags().StringSliceVarP(&vars.contexts, contextFlag, contextFlagShort, nil, destroyContextDescription)
	cmd.Flags().BoolVar(&vars.destroyForce, destroyContextForceFlag, false, destroyContextForceDescription)
	_ = cmd.RegisterFlagCompletionFunc(contextFlag, NewContextAutoComplete().GetContextAutoComplete())
	return cmd
}
