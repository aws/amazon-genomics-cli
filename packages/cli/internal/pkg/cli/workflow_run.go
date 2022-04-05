// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	contextFlag            = "context"
	contextFlagShort       = "c"
	contextFlagDescription = "Name of context"
)

const (
	inputsFileFlag            = "inputsFile"
	inputsFileFlagShort       = "i"
	inputsFileFlagDescription = "Inputs File Path"
)

type runWorkflowVars struct {
	WorkflowName string
	Arguments    string
	OptionFile   string
	ContextName  string
}

type runWorkflowOpts struct {
	runWorkflowVars
	wfManager *workflow.Manager
}

func newRunWorkflowOpts(vars runWorkflowVars) (*runWorkflowOpts, error) {
	return &runWorkflowOpts{
		runWorkflowVars: vars,
		wfManager:       workflow.NewManager(profile),
	}, nil
}

func (o *runWorkflowOpts) Validate() error {
	return nil
}

func (o *runWorkflowOpts) Execute() (string, error) {
	return o.wfManager.RunWorkflow(o.ContextName, o.WorkflowName, o.Arguments, o.OptionFile)
}

func BuildWorkflowRunCommand() *cobra.Command {
	vars := runWorkflowVars{}
	cmd := &cobra.Command{
		Use:   "run workflow_name --context context_name",
		Short: "Run a workflow",
		Long: `run is for running the specified workflow in the specified context.
This command prints a run Id for the created workflow instance.
`,
		Example: `
Run the workflow named "example-workflow", against the "prod" context,
using input parameters contained in file "file:///Users/ec2-user/myproj/test-args.json"
/code $ agc workflow run example-workflow --context prod --inputsFile file:///Users/ec2-user/myproj/test-args.json`,
		Args: cobra.ExactArgs(1),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			vars.WorkflowName = args[0]
			opts, err := newRunWorkflowOpts(vars)
			if err != nil {
				return clierror.New("workflow run", vars, err)
			}
			log.Info().Msgf("Running workflow. Workflow name: '%s', InputsFile: '%s', OptionFile: '%s', Context: '%s'",
				opts.WorkflowName, opts.Arguments, opts.OptionFile, opts.ContextName)
			if err := opts.Validate(); err != nil {
				return err
			}
			instanceId, err := opts.Execute()
			if err != nil {
				return clierror.New("workflow run", vars, err)
			}
			format.Default.Write(instanceId)
			return nil
		}),
		ValidArgsFunction: NewWorkflowAutoComplete().GetWorkflowAutoComplete(),
	}
	cmd.Flags().StringVarP(&vars.Arguments, inputsFileFlag, inputsFileFlagShort, "", inputsFileFlagDescription)
	cmd.Flags().StringVarP(&vars.OptionFile, optionFileFlag, optionFileFlagShort, "", optionFileFlagDescription)
	cmd.Flags().StringVarP(&vars.ContextName, contextFlag, contextFlagShort, "", contextFlagDescription)
	aliasFn := func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		switch name {
		case argsFlag:
			name = inputsFileFlag
		}
		return pflag.NormalizedName(name)
	}
	cmd.Flags().SetNormalizeFunc(aliasFn)
	_ = cmd.MarkFlagRequired(contextFlag)
	_ = cmd.RegisterFlagCompletionFunc(contextFlag, NewContextAutoComplete().GetContextAutoComplete())
	return cmd
}
