// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
	"github.com/aws/amazon-genomics-cli/internal/pkg/stringutils"
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
	engine        string
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

	ctxMap, err := o.ctxManager.List()
	if err != nil {
		return err
	}

	summary := ctxMap[o.contextName]
	o.engine = summary.Engines[0].Engine

	if o.workflowRunId == "" {
		if !summary.IsServerProcessEngine() {
			log.Warn().Msgf("DEPRECATION WARNING!!")
			log.Warn().Msgf("Obtaining engine logs for a workflow context where the engine is '%s' will return engine logs from ALL workflows run in this context",
				o.engine)
			log.Warn().Msgf("Specifying a run-id will be MANDATORY in future versions")
			log.Warn().Msgf("Please run the command again with -r <run-id>")
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

	if o.workflowRunId == "" {
		return executeGetEngineLogForWholeGroup(o, logGroupName)
	}
	if o.engine == constants.CROMWELL {
		currentFilter := o.filter
		log.Info().Msgf("current filter = '%s'", currentFilter)
		cromwellLogTag := stringutils.SubString(o.workflowRunId, 0, 7)
		log.Info().Msgf("cromwell short run id = '%s'", cromwellLogTag)

		var builder strings.Builder
		builder.WriteString(currentFilter)
		builder.WriteByte(' ')
		builder.WriteString(cromwellLogTag)
		o.filter = builder.String()
		log.Info().Msgf("filter = '%s'", o.filter)

		return executeGetEngineLogForWholeGroup(o, logGroupName)
	}
	return executeGetEngineLogForRunId(o, logGroupName)
}

func executeGetEngineLogForWholeGroup(o *logsEngineOpts, logGroupName string) error {
	if o.tail {
		return o.followLogGroup(logGroupName)
	}
	return o.displayLogGroup(logGroupName, o.startTime, o.endTime, o.filter)
}

func executeGetEngineLogForRunId(o *logsEngineOpts, logGroupName string) error {
	log.Info().Msgf("Getting log stream for workflow run '%s'", o.workflowRunId)

	workflowRunLog, err := o.workflowManager.GetEngineLogByRunId(o.workflowRunId)
	if err != nil {
		return err
	}

	logStreamNames := streamNamesFromRunLog(workflowRunLog)
	log.Debug().Msgf("Log stream name is: '%v'", logStreamNames)

	if len(logStreamNames) > 0 {
		if o.tail {
			return o.followLogStreams(logGroupName, logStreamNames...)
		}
		return o.displayLogStreams(logGroupName, o.startTime, o.endTime, o.filter, logStreamNames...)
	}

	workflowStatus := workflowRunLog.WorkflowStatus
	log.Warn().Msgf("Cannot find an engine log stream for workflow run '%s', the current status of the run is: '%s', "+
		"the log will not be available until after the workflow is RUNNING", o.workflowRunId, workflowStatus)
	return nil
}

func streamNamesFromRunLog(workflowRunLog workflow.EngineLog) []string {
	var streamNames = make([]string, 0)
	if workflowRunLog.StdOut != "" {
		streamNames = append(streamNames, workflowRunLog.StdOut)
	}
	if workflowRunLog.StdErr != "" {
		streamNames = append(streamNames, workflowRunLog.StdErr)
	}
	return streamNames
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
