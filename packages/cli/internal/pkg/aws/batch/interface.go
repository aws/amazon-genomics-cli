package batch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/batch"
)

type Interface interface {
	GetJobs(jobIds []string) ([]Job, error)
}

type batchInterface interface {
	DescribeJobs(context.Context, *batch.DescribeJobsInput, ...func(*batch.Options)) (*batch.DescribeJobsOutput, error)
}
