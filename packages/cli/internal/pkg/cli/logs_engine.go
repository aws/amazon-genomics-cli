// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/common/aws"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type logsEngineVars struct {
	logsSharedVars
}

type logsEngineOpts struct {
	logsEngineVars
	logsSharedOpts
}

func newLogsEngineOpts(vars logsEngineVars) (*logsEngineOpts, error) {
	return &logsEngineOpts{
		logsEngineVars: vars,
		logsSharedOpts: logsSharedOpts{
			ctxManager: context.NewManager(profile),
			cwlClient:  aws.CwlClient(profile),
		},
	}, nil
}

func (o *logsEngineOpts) Validate() error {
	if err := o.validateFlags(); err != nil {
		return err
	}

	return o.parseTime(o.logsSharedVars)
}

func (o *logsEngineOpts) Execute() error {
	contextInfo, err := o.ctxManager.Info(o.contextName)
	if err != nil {
		return err
	}

	logGroupName := contextInfo.EngineLogGroupName
	if o.tail {
		err = o.followLogGroup(logGroupName)
	} else {
		err = o.displayLogGroup(logGroupName, o.startTime, o.endTime, o.filter)
	}

	return err
}

func BuildLogsEngineCommand() *cobra.Command {
	vars := logsEngineVars{}
	cmd := &cobra.Command{
		Use:   "engine -c context_name [-f filter] [-s start_date] [-e end_date] [-l look_back] [-t]",
		Short: "Show workflow engine logs for a given context.",
		Long: `Show workflow engine logs for a given context.
If no start, end, or look back periods are set, this command will show logs from the last hour.`,
		Example: `
/code agc logs engine -c myCtx -s 2021/3/31 -e 2021/4/1 -f ERROR`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newLogsEngineOpts(vars)
			if err != nil {
				return err
			}
			if err = opts.Validate(); err != nil {
				return err
			}
			opts.setDefaultEndTimeIfEmpty()
			log.Info().Msgf("Showing engine logs for '%s'", opts.contextName)
			if err = opts.Execute(); err != nil {
				return clierror.New("logs engine", vars, err, "check you have valid aws credentials, check that the named context is defined in the agc-project.yaml file, check the context is deployed")
			}
			return nil
		}),
	}
	vars.setFilterFlags(cmd)
	vars.setContextFlag(cmd)
	return cmd
}
