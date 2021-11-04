package cwl

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/util"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type GetLogsInput struct {
	LogGroupName string
	StartTime    *time.Time
	EndTime      *time.Time
	Filter       string
	Streams      []string
}

type GetLogsOutput struct {
	paginator *cloudwatchlogs.FilterLogEventsPaginator
}

func (o GetLogsOutput) HasMoreLogs() bool {
	return o.paginator.HasMorePages()
}

func (o GetLogsOutput) NextLogs() ([]string, error) {
	var logs []string
	output, err := o.paginator.NextPage(context.Background())
	if err != nil {
		return nil, actionableerror.FindSuggestionForError(err, actionableerror.AwsErrorMessageToSuggestedActionMap)
	}
	logs = append(logs, formatEvents(output.Events)...)
	return logs, nil
}

func (c Client) GetLogsPaginated(input GetLogsInput) LogPaginator {
	eventInput := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:   aws.String(input.LogGroupName),
		StartTime:      util.TimeToAws(input.StartTime),
		EndTime:        util.TimeToAws(input.EndTime),
		FilterPattern:  aws.String(input.Filter),
		LogStreamNames: input.Streams,
	}
	paginator := cloudwatchlogs.NewFilterLogEventsPaginator(c.cwl, eventInput)
	return GetLogsOutput{paginator: paginator}
}

func formatEvents(events []types.FilteredLogEvent) []string {
	logsByStream := make(map[string][]*types.FilteredLogEvent)
	for index := range events {
		event := events[index]
		logsByStream[*event.LogStreamName] = append(logsByStream[*event.LogStreamName], &event)
	}

	return convertStreamLogsToLogs(logsByStream, len(events))
}

func convertStreamLogsToLogs(logsByStream map[string][]*types.FilteredLogEvent, eventSize int) []string {
	logs, index := make([]string, eventSize), 0
	for _, eventList := range logsByStream {
		for _, event := range eventList {
			logs[index] = formatEvent(*event)
			index++
		}
	}
	return logs
}

func formatEvent(event types.FilteredLogEvent) string {
	timestamp := util.TimeFromAws(event.Timestamp)
	return fmt.Sprintf("%s\t%s", timestamp.Format(time.RFC1123Z), aws.ToString(event.Message))
}
