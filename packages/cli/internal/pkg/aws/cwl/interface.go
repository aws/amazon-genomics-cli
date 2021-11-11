package cwl

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type LogPaginator interface {
	HasMoreLogs() bool
	NextLogs() ([]string, error)
}

type Interface interface {
	GetLogsPaginated(input GetLogsInput) LogPaginator
	StreamLogs(ctx context.Context, logGroupName string, streams ...string) <-chan StreamEvent
}

type cwlInterface interface {
	cloudwatchlogs.FilterLogEventsAPIClient
}
