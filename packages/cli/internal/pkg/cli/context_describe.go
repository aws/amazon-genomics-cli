// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/format"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type describeContextVars struct {
	ContextName string
}

type describeContextOpts struct {
	describeContextVars
	ctxManager context.Interface
}

func newDescribeContextOpts(vars describeContextVars) (*describeContextOpts, error) {
	return &describeContextOpts{
		describeContextVars: vars,
		ctxManager:          context.NewManager(profile),
	}, nil
}

func (o *describeContextOpts) Validate(args []string) error {
	if len(args) == 0 && o.ContextName == "" {
		return fmt.Errorf("a context must be provided")
	} else if len(args) == 1 && o.ContextName != "" {
		return fmt.Errorf("either the '-c' flag or a context must be provided, but not both")
	}

	if len(args) == 1 {
		o.ContextName = args[0]
	}

	return nil
}

// Execute returns a context definition for the specified name.
func (o *describeContextOpts) Execute() (types.Context, error) {
	ctxName := o.ContextName
	info, err := o.ctxManager.Info(ctxName)
	if err != nil {
		return types.Context{}, err
	}
	return types.Context{
		Name:                 ctxName,
		Status:               info.Status.ToString(),
		StatusReason:         info.StatusReason,
		InstanceTypes:        buildInstanceTypes(info.InstanceTypes),
		MaxVCpus:             info.MaxVCpus,
		RequestSpotInstances: info.IsSpot,
		Output:               types.OutputLocation{Url: info.BucketLocation},
	}, nil
}

func buildInstanceTypes(stringTypes []string) []types.InstanceType {
	var instanceTypes []types.InstanceType
	for _, val := range stringTypes {
		instanceTypes = append(instanceTypes, types.InstanceType{Value: val})
	}
	return instanceTypes
}

// BuildContextDescribeCommand builds the command to show the information for a specific or for all contexts in the current project.
func BuildContextDescribeCommand() *cobra.Command {
	vars := describeContextVars{}
	cmd := &cobra.Command{
		Use:   "describe context_name",
		Short: "Show the information for a specific context in the current project",
		Long: `describe is for showing information about the specified context.

` + DescribeOutput(types.Context{}),
		Example: `
/code agc context describe myCtx`,
		Args: cobra.ArbitraryArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newDescribeContextOpts(vars)
			if err != nil {
				return err
			}
			if err := opts.Validate(args); err != nil {
				return err
			}
			log.Info().Msgf("Describing context '%s'", opts.ContextName)
			ctx, err := opts.Execute()
			if err != nil {
				return clierror.New("context describe", vars, err)
			}
			format.Default.Write(ctx)
			return nil
		}),
	}
	cmd.Flags().StringVarP(&vars.ContextName, contextFlag, contextFlagShort, "", deployContextDescription)
	return cmd
}
