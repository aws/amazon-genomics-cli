package batch

import (
	"context"
	"fmt"
	"testing"
	"time"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
	"github.com/stretchr/testify/assert"
)

var (
	testJobId           = "test-job-id"
	testJobName         = "test-job-name"
	testJobStatus       = "test-job-status"
	testJobStatusReason = "test-job-status-reason"
	testJobCommand      = "test-job-command"
	testJobStreamName   = "test-job-stream-name"
	testStartTime       = time.Now().Add(-time.Hour)
	testStopTime        = time.Now()
)

func (m *BatchMock) DescribeJobs(ctx context.Context, input *batch.DescribeJobsInput, _ ...func(*batch.Options)) (*batch.DescribeJobsOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*batch.DescribeJobsOutput), err
	}
	return nil, err
}

func TestClient_GetJobs(t *testing.T) {
	client := NewMockClient()
	client.batch.(*BatchMock).On("DescribeJobs", context.Background(), &batch.DescribeJobsInput{
		Jobs: []string{testJobId},
	}).Return(&batch.DescribeJobsOutput{
		Jobs: []types.JobDetail{{
			JobId:     &testJobId,
			JobName:   &testJobName,
			StartedAt: testStartTime.UnixNano() / 1000000,
			Status:    types.JobStatus(testJobStatus),
			Container: &types.ContainerDetail{
				Command:       []string{testJobCommand},
				LogStreamName: &testJobStreamName,
			},
			StatusReason: &testJobStatusReason,
			StoppedAt:    testStopTime.UnixNano() / 1000000,
		}},
	}, nil)

	jobs, err := client.GetJobs([]string{testJobId})

	actualStartTime := jobs[0].StartTime
	actualStopTime := jobs[0].StopTime
	jobs[0].StartTime = nil
	jobs[0].StopTime = nil

	assert.NoError(t, err)
	assert.True(t, actualStartTime.Equal(testStartTime.Truncate(time.Millisecond)))
	assert.True(t, actualStopTime.Equal(testStopTime.Truncate(time.Millisecond)))
	assert.Equal(t, []Job{{
		JobId:         testJobId,
		JobName:       testJobName,
		Commands:      []string{testJobCommand},
		StartTime:     nil,
		StopTime:      nil,
		JobStatus:     types.JobStatus(testJobStatus),
		StatusReason:  testJobStatusReason,
		LogStreamName: testJobStreamName,
	}}, jobs)
}

func TestClient_GetJobs_Error(t *testing.T) {
	client := NewMockClient()
	client.batch.(*BatchMock).On("DescribeJobs", context.Background(), &batch.DescribeJobsInput{
		Jobs: []string{testJobId},
	}).Return(nil, fmt.Errorf("some job error"))

	jobs, err := client.GetJobs([]string{testJobId})

	assert.Equal(t, fmt.Errorf("some job error"), err)
	assert.Empty(t, jobs)
}
