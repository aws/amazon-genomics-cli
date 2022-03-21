// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	runIdFlag        = "run-id"
	runIdShort       = "r"
	runIdDescription = "filter to engine logs to this workflow run id"
	runIdDefault     = ""
)

type logsEngineVars struct {
	logsSharedVars
	workflowRunId string
}

type logsEngineOpts struct {
	logsEngineVars
	logsSharedOpts
	workflowManager *workflow.Manager
}

func newLogsEngineOpts(vars logsEngineVars) (*logsEngineOpts, error) {
	return &logsEngineOpts{
		logsEngineVars: vars,
		logsSharedOpts: logsSharedOpts{
			ctxManager: context.NewManager(profile),
			cwlClient:  aws.CwlClient(profile),
		},
		workflowManager: workflow.NewManager(profile),
	}, nil
}

func (o *logsEngineOpts) Validate() error {
	if err := o.validateFlags(); err != nil {
		return err
	}

	if o.workflowRunId == "" {
		ctxMap, err := o.ctxManager.List()
		if err != nil {
			return err
		}

		summary := ctxMap[o.contextName]
		if summary.IsHeadProcessEngine() {
			return actionableerror.New(fmt.Errorf("a workflow run must be specified if workflow engine is '%s'", summary.Engines[0].Engine),
				"please run the command again with -r <run-id>")
		}
	}

	return o.parseTime(o.logsSharedVars)
}

func (o *logsEngineOpts) Execute() error {
	contextInfo, err := o.ctxManager.Info(o.contextName)
	if err != nil {
		return err
	}

	logGroupName := contextInfo.EngineLogGroupName
	log.Debug().Msgf("Engine log group name: '%s'", logGroupName)

	if o.workflowRunId != "" {
		log.Debug().Msgf("Getting log stream for workflow run '%s'", o.workflowRunId)

		workflowRunLog, err := o.workflowManager.GetEngineLogByRunId(o.workflowRunId)
		if err != nil {
			return err
		}
		log.Debug().Msgf("Stream '%s' for log group '%s' contains StdOut for run '%s'", workflowRunLog.StdOut, logGroupName, workflowRunLog.WorkflowRunId)

		if o.tail {
			err = o.followLogStreams(logGroupName, workflowRunLog.StdOut)
		} else {
			err = o.displayLogStreams(logGroupName, o.startTime, o.endTime, o.filter, workflowRunLog.StdOut)
		}

	} else {
		if o.tail {
			err = o.followLogGroup(logGroupName)
		} else {
			err = o.displayLogGroup(logGroupName, o.startTime, o.endTime, o.filter)
		}
	}

	return err
}

func BuildLogsEngineCommand() *cobra.Command {
	vars := logsEngineVars{}
	cmd := &cobra.Command{
		Use:   "engine -c context_name [-r run_id] [-f filter] [-s start_date] [-e end_date] [-l look_back] [-t]",
		Short: "Show workflow engine logs for a given context.",
		Long: `Show workflow engine logs for a given context.
If no start, end, or look back periods are set, this command will show logs from the last hour.`,
		Example: `
/code agc logs engine -c myCtx -r 1234-aed-32db -s 2021/3/31 -e 2021/4/1 -f ERROR`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newLogsEngineOpts(vars)
			if err != nil {
				return err
			}
			if err = opts.Validate(); err != nil {
				return err
			}

			if vars.workflowRunId == "" {
				opts.setDefaultEndTimeIfEmpty()
			}
			var msg = fmt.Sprintf("Showing engine logs for '%s'", opts.contextName)
			if vars.endString != "" {
				msg = fmt.Sprintf("Showing engine logs for '%s' from '%s", opts.contextName, opts.endString)
			}
			log.Info().Msg(msg)
			if err = opts.Execute(); err != nil {
				return clierror.New("logs engine", vars, err)
			}
			return nil
		}),
	}
	vars.setFilterFlags(cmd)
	vars.setContextFlag(cmd)
	cmd.Flags().StringVarP(&vars.workflowRunId, runIdFlag, runIdShort, runIdDefault, runIdDescription)
	return cmd
}
