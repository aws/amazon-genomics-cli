package cwl

import (
	"context"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/util"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

var sleepDuration = time.Second * 1

type StreamEvent struct {
	Logs []string
	Err  error
}

func (c Client) StreamLogs(ctx context.Context, logGroupName string, streams ...string) <-chan StreamEvent {
	stream := make(chan StreamEvent)
	go func() {
		defer func() { close(stream) }()
		now := time.Now()
		startTime := util.TimeToAws(&now)
		var lastToken *string
		for {
			output, err := c.cwl.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName:   aws.String(logGroupName),
				StartTime:      startTime,
				LogStreamNames: streams,
			})
			if err != nil {
				stream <- StreamEvent{Err: actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)}
				return
			}
			if aws.ToString(lastToken) != aws.ToString(output.NextToken) {
				stream <- StreamEvent{Logs: parseEventLogs(output.Events, startTime)}
				lastToken = output.NextToken
			}

			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(sleepDuration)
			}
		}
	}()
	return stream
}

func parseEventLogs(events []types.FilteredLogEvent, latestTimestamp *int64) []string {
	logs := make([]string, len(events))
	for i, event := range events {
		if aws.ToInt64(event.Timestamp) > aws.ToInt64(latestTimestamp) {
			latestTimestamp = event.Timestamp
		}
		logs[i] = formatEvent(event)
	}
	return logs
}
