package batch

import (
	"context"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/util"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
)

type Job struct {
	JobId         string
	JobName       string
	Commands      []string
	StartTime     *time.Time
	StopTime      *time.Time
	JobStatus     types.JobStatus
	StatusReason  string
	LogStreamName string
}

func (c Client) GetJobs(jobIds []string) ([]Job, error) {
	input := &batch.DescribeJobsInput{Jobs: jobIds}
	output, err := c.batch.DescribeJobs(context.Background(), input)
	if err != nil {
		return nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	jobs := make([]Job, len(output.Jobs))
	for i, job := range output.Jobs {
		startTime := util.TimeFromAws(&job.StartedAt)
		stopTime := util.TimeFromAws(&job.StoppedAt)
		jobs[i] = Job{
			JobId:         aws.ToString(job.JobId),
			JobName:       aws.ToString(job.JobName),
			Commands:      job.Container.Command,
			StartTime:     &startTime,
			StopTime:      &stopTime,
			JobStatus:     job.Status,
			StatusReason:  aws.ToString(job.StatusReason),
			LogStreamName: aws.ToString(job.Container.LogStreamName),
		}
	}
	return jobs, nil
}
