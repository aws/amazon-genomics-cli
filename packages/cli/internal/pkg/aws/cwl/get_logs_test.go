package cwl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/assert"
)

const (
	testLogGroupName  = "test-log-group-name"
	testFilterPattern = "test-filter-pattern"
)

func (m *CwlMock) FilterLogEvents(ctx context.Context, input *cloudwatchlogs.FilterLogEventsInput, _ ...func(options *cloudwatchlogs.Options)) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*cloudwatchlogs.FilterLogEventsOutput), err
	}
	return nil, err
}

func TestClient_GetLogs(t *testing.T) {
	client := NewMockClient()
	endTime := time.Now()
	startTime := endTime.Add(-time.Second)
	eventTime1 := time.Unix(0, 773391600000000000)
	eventTime2 := time.Unix(0, 773392000000000000)
	client.cwl.(*CwlMock).On("FilterLogEvents", context.Background(), &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  aws.String(testLogGroupName),
		StartTime:     aws.Int64(startTime.UnixNano() / 1000000),
		EndTime:       aws.Int64(endTime.UnixNano() / 1000000),
		FilterPattern: aws.String(testFilterPattern),
	}).Return(&cloudwatchlogs.FilterLogEventsOutput{Events: []types.FilteredLogEvent{
		{
			EventId:       aws.String("some-id"),
			IngestionTime: aws.Int64(eventTime1.UnixNano() / 1000000),
			LogStreamName: aws.String("log-stream-1"),
			Message:       aws.String("Hello"),
			Timestamp:     aws.Int64(eventTime1.UnixNano() / 1000000),
		},
		{
			EventId:       aws.String("some-other-id"),
			IngestionTime: aws.Int64(eventTime2.UnixNano() / 1000000),
			LogStreamName: aws.String("log-stream-1"),
			Message:       aws.String("world!"),
			Timestamp:     aws.Int64(eventTime2.UnixNano() / 1000000),
		},
	}}, nil)
	output := client.GetLogsPaginated(GetLogsInput{
		LogGroupName: testLogGroupName,
		StartTime:    &startTime,
		EndTime:      &endTime,
		Filter:       testFilterPattern,
	})
	assert.True(t, output.HasMoreLogs())
	logs, err := output.NextLogs()
	assert.False(t, output.HasMoreLogs())
	assert.NoError(t, err)
	assert.Equal(t, []string{fmt.Sprintf("%s\tHello", eventTime1.Format(time.RFC1123Z)), fmt.Sprintf("%s\tworld!", eventTime2.Format(time.RFC1123Z))}, logs)
}

func TestClient_GetLogs_Error(t *testing.T) {
	client := NewMockClient()
	client.cwl.(*CwlMock).On("FilterLogEvents", context.Background(), &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  aws.String(testLogGroupName),
		FilterPattern: aws.String(testFilterPattern),
	}).Return(nil, fmt.Errorf("some log error"))
	output := client.GetLogsPaginated(GetLogsInput{
		LogGroupName: testLogGroupName,
		StartTime:    nil,
		EndTime:      nil,
		Filter:       testFilterPattern,
	})
	logs, err := output.NextLogs()
	assert.Equal(t, fmt.Errorf("some log error"), err)
	assert.Empty(t, logs)
}
