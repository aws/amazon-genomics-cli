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
	"github.com/stretchr/testify/mock"
)

func TestClient_StreamLogs(t *testing.T) {
	someTime := time.Unix(0, 773391600000000000)
	ctx, cancel := context.WithCancel(context.Background())
	client := NewMockClient()
	client.cwl.(*CwlMock).On("FilterLogEvents", ctx, mock.Anything).
		Return(&cloudwatchlogs.FilterLogEventsOutput{
			NextToken: aws.String("Token"),
			Events: []types.FilteredLogEvent{
				{
					EventId:       aws.String("ID"),
					IngestionTime: aws.Int64(someTime.UnixNano() / 1000000),
					LogStreamName: aws.String("log-stream-1"),
					Message:       aws.String("Hello!"),
					Timestamp:     aws.Int64(someTime.UnixNano() / 1000000),
				},
			}}, nil)
	cancel()
	stream := client.StreamLogs(ctx, testLogGroupName)
	event := <-stream
	assert.Equal(t, []string{fmt.Sprintf("%s\tHello!", someTime.Format(time.RFC1123Z))}, event.Logs)
	assert.NoError(t, event.Err)
	cancel()
	_, isOpen := <-stream
	assert.False(t, isOpen)
}
func TestClient_StreamLogs_EmptyLog(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client := NewMockClient()
	client.cwl.(*CwlMock).On("FilterLogEvents", ctx, mock.Anything).
		Return(&cloudwatchlogs.FilterLogEventsOutput{
			NextToken: aws.String("Token"),
			Events:    []types.FilteredLogEvent{}}, nil)
	cancel()
	stream := client.StreamLogs(ctx, testLogGroupName)
	event := <-stream
	assert.Equal(t, []string{}, event.Logs)
	assert.NoError(t, event.Err)
	cancel()
	_, isOpen := <-stream
	assert.False(t, isOpen)
}

func TestClient_StreamLogs_Error(t *testing.T) {
	ctx := context.Background()
	client := NewMockClient()
	client.cwl.(*CwlMock).On("FilterLogEvents", ctx, mock.Anything).
		Return(nil, fmt.Errorf("some log error"))
	stream := client.StreamLogs(ctx, testLogGroupName)
	event := <-stream
	assert.Equal(t, fmt.Errorf("some log error"), event.Err)
	_, isOpen := <-stream
	assert.False(t, isOpen)
}
