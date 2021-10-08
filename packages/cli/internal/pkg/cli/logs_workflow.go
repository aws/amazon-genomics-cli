// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/batch"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	logWorkflowRunFlag            = "run"
	logWorkflowRunFlagShort       = "r"
	logWorkflowRunFlagDescription = `The ID of a workflow run to retrieve.`

	logFailedTasksFlag            = "failed-tasks"
	logFailedTasksFlagDescription = `Only show logs of tasks that have not exited cleanly.`

	cachedJobId = "XXXXX"
)

var noRunsFoundError = errors.New("there are no runs for the workflow")

type logsWorkflowVars struct {
	logsSharedVars
	workflowName string
	runId        string
	failedTasks  bool
}

type logsWorkflowOpts struct {
	logsWorkflowVars
	logsSharedOpts
	batchClient     batch.Interface
	workflowManager workflow.TasksManager
}

func newLogsWorkflowOpts(vars logsWorkflowVars) (*logsWorkflowOpts, error) {
	return &logsWorkflowOpts{
		logsWorkflowVars: vars,
		logsSharedOpts: logsSharedOpts{
			ctxManager: context.NewManager(profile),
			cwlClient:  aws.CwlClient(profile),
		},
		batchClient:     aws.BatchClient(profile),
		workflowManager: workflow.NewManager(profile),
	}, nil
}

func (o *logsWorkflowOpts) Validate() error {
	if err := o.validateFlags(); err != nil {
		return err
	}

	return o.parseTime(o.logsSharedVars)
}

func (o *logsWorkflowOpts) Execute() error {
	err := o.setRunId()
	if err != nil {
		if errors.Is(err, noRunsFoundError) {
			log.Info().Msgf("Workflow '%s' has not been run yet", o.workflowName)
			return nil
		}
		return err
	}
	log.Debug().Msgf("Showing logs for workflow run '%s'", o.runId)
	jobIds, err := o.getJobIds()
	if err != nil {
		return err
	}

	if len(jobIds) == 0 {
		log.Info().Msgf("No logs available for run '%s'. Please try again later.", o.runId)
		return nil
	}
	notCachedJobIds := filterCachedJobIds(jobIds)
	totalJobs := len(jobIds)
	notCachedJobs := len(notCachedJobIds)
	cachedJobs := totalJobs - notCachedJobs
	if cachedJobs > 0 {
		log.Info().Msgf("%d of %d jobs were cached. Logs are not available for cached jobs", cachedJobs, totalJobs)
	}
	if notCachedJobs == 0 {
		return nil
	}

	streams, err := o.getStreamsForJobs(notCachedJobIds)
	if err != nil {
		return err
	}

	logGroupName := "/aws/batch/job"
	if o.tail {
		_ = o.followLogGroup(logGroupName, streams...)
	} else {
		_ = o.displayLogGroup(logGroupName, o.startTime, o.endTime, o.filter, streams...)
	}

	return nil
}

func filterCachedJobIds(ids []string) []string {
	var result []string
	for _, id := range ids {
		if id != cachedJobId {
			result = append(result, id)
		}
	}
	return result
}

func (o *logsWorkflowOpts) setRunId() error {
	if o.runId == "" {
		instances, err := o.workflowManager.StatusWorkflowByName(o.workflowName, 1)
		if err != nil {
			return err
		}
		if len(instances) == 0 {
			return noRunsFoundError
		}
		o.runId = instances[0].Id
		log.Info().Msgf("Showing logs for the latest run of the workflow. Run id: '%s'", o.runId)
	}
	return nil
}

func (o *logsWorkflowOpts) getJobIds() ([]string, error) {
	tasks, err := o.workflowManager.GetWorkflowTasks(o.runId)
	if err != nil {
		return nil, err
	}

	var jobIds []string
	for _, task := range tasks {
		if o.failedTasks && task.ExitCode == 0 {
			log.Debug().Msgf("skipping successful task '%s' ('%s')", task.Name, task.JobId)
			continue
		}
		if task.JobId == cachedJobId {
			log.Debug().Msgf("skipping cached task '%s'", task.Name)
		}
		jobIds = append(jobIds, task.JobId)
	}
	return jobIds, nil
}

func (o *logsWorkflowOpts) getStreamsForJobs(jobIds []string) ([]string, error) {
	jobs, err := o.batchClient.GetJobs(jobIds)
	if err != nil {
		return nil, err
	}
	streams := make([]string, len(jobs))
	for i, job := range jobs {
		if job.LogStreamName == "" {
			log.Debug().Msgf("No log stream found for job '%s' ('%s')", job.JobName, job.JobId)
			continue
		}
		streams[i] = job.LogStreamName
	}
	return streams, err
}

// BuildLogsWorkflowCommand builds the command to output the content of Cloudwatch log streams
// of workflows.
func BuildLogsWorkflowCommand() *cobra.Command {
	vars := logsWorkflowVars{}
	cmd := &cobra.Command{
		Use:   "workflow workflow_name [-r run_id] [--failed_tasks]",
		Short: "Show the task logs of a given workflow",
		Long: `Show the task logs of a given workflow.
If the --run flag is omitted then the latest workflow run is used.`,
		Args: cobra.ExactArgs(1),
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			vars.workflowName = args[0]
			opts, err := newLogsWorkflowOpts(vars)
			if err != nil {
				return err
			}
			if err = opts.Validate(); err != nil {
				return err
			}
			log.Info().Msgf("Showing the logs for '%s'", vars.workflowName)
			if err = opts.Execute(); err != nil {
				return clierror.New("logs workflow", vars, err)
			}
			return nil
		}),
		ValidArgsFunction: NewWorkflowAutoComplete().GetWorkflowAutoComplete(),
	}
	vars.setFilterFlags(cmd)
	cmd.Flags().StringVarP(&vars.runId, logWorkflowRunFlag, logWorkflowRunFlagShort, "", logWorkflowRunFlagDescription)
	cmd.Flags().BoolVar(&vars.failedTasks, logFailedTasksFlag, false, logFailedTasksFlagDescription)
	return cmd
}
