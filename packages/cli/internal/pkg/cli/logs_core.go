package cli

import (
	ctx "context"
	"fmt"
	"sync"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	logStartFlag            = "start"
	logStartFlagShort       = "s"
	logStartFlagDescription = `A date to begin displaying logs.
Supports most date formats, such as 2021/03/31 or 8/8/2021 01:00:01 PM.
Times respect the system timezone.`

	logEndFlag            = "end"
	logEndFlagShort       = "e"
	logEndFlagDescription = `A date to stop displaying logs.
Supports most date formats, such as mm/dd/yy or yyyy-mm-dd-07:00.
Times respect the system timezone.`

	logLookBackFlag            = "look-back"
	logLookBackFlagShort       = "l"
	logLookBackFlagDescription = `A period of time to look back from now, such as "2h45m".
Valid time units are "s", "m", and "h".`

	logFilterFlag            = "filter"
	logFilterFlagShort       = "f"
	logFilterFlagDescription = `Match terms, phrases, or values in the logs.
Filters are case sensitive and multiple terms combine with AND logic.
Use a question mark for OR, such as "?ERROR ?WARN". Filter out terms with a minus, such as "-INFO".`

	tailFlag            = "tail"
	tailFlagShort       = "t"
	tailFlagDescription = "Follow the log output."
)

var logInfo = log.Info

var printLn = func(args ...interface{}) {
	_, _ = fmt.Println(args...)
}

type logsSharedVars struct {
	tail        bool
	contextName string
	startString string
	endString   string
	lookBack    string
	filter      string
}

type logsSharedOpts struct {
	startTime  *time.Time
	endTime    *time.Time
	ctxManager context.Interface
	cwlClient  cwl.Interface
}

var now = time.Now

func (v *logsSharedVars) validateFlags() error {
	if (v.startString != "" || v.endString != "") && v.lookBack != "" {
		return fmt.Errorf("a look back period cannot be specified together with start or end times")
	}
	return nil
}

func (v *logsSharedVars) setFilterFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&v.tail, tailFlag, tailFlagShort, false, tailFlagDescription)
	cmd.Flags().StringVarP(&v.startString, logStartFlag, logStartFlagShort, "", logStartFlagDescription)
	cmd.Flags().StringVarP(&v.endString, logEndFlag, logEndFlagShort, "", logEndFlagDescription)
	cmd.Flags().StringVarP(&v.lookBack, logLookBackFlag, logLookBackFlagShort, "", logLookBackFlagDescription)
	cmd.Flags().StringVarP(&v.filter, logFilterFlag, logFilterFlagShort, "", logFilterFlagDescription)
}

func (v *logsSharedVars) setContextFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&v.contextName, contextFlag, contextFlagShort, "", contextFlagDescription)
	_ = cmd.MarkFlagRequired(contextFlag)
	_ = cmd.RegisterFlagCompletionFunc(contextFlag, NewContextAutoComplete().GetContextAutoComplete())
}

func (o *logsSharedOpts) setDefaultEndTimeIfEmpty() {
	if o.startTime == nil && o.endTime == nil {
		lastHour := now().Add(-1 * time.Hour)
		o.startTime = &lastHour
	}
}

func (o *logsSharedOpts) parseTime(vars logsSharedVars) error {
	if vars.startString != "" {
		t, err := dateparse.ParseLocal(vars.startString)
		if err != nil {
			return err
		}
		o.startTime = &t
	}
	if vars.endString != "" {
		t, err := dateparse.ParseLocal(vars.endString)
		if err != nil {
			return err
		}
		o.endTime = &t
	}
	if vars.lookBack != "" {
		lookBack, err := time.ParseDuration(vars.lookBack)
		if err != nil {
			return err
		}
		then := time.Now().Add(-lookBack)
		o.startTime = &then
	}
	return nil
}

func (o *logsSharedOpts) followLogGroup(logGroupName string) error {
	channel := o.cwlClient.StreamLogs(ctx.Background(), logGroupName)
	return o.displayEventFromChannel(channel)
}

func (o *logsSharedOpts) displayEventFromChannel(channel <-chan cwl.StreamEvent) error {
	firstEvent := true

	for event := range channel {
		if event.Err != nil {
			return event.Err
		}
		if len(event.Logs) > 0 {
			for _, line := range event.Logs {
				printLn(line)
			}
		} else if firstEvent {
			logInfo().Msg("There are no new logs. Please wait for the first logs to appear...")
		} else {
			log.Debug().Msg("No new logs")
		}

		firstEvent = false
	}

	return nil
}

func (o *logsSharedOpts) followLogStreams(logGroupName string, streams ...string) error {
	const maxLogStreams = 100
	streamingCtx, cancelFunc := ctx.WithCancel(ctx.Background())
	defer cancelFunc()
	streamBatches := splitToBatchesBy(maxLogStreams, streams)
	var eventChannels []<-chan cwl.StreamEvent

	for _, batch := range streamBatches {
		eventChannel := o.cwlClient.StreamLogs(streamingCtx, logGroupName, batch...)
		eventChannels = append(eventChannels, eventChannel)
	}

	return o.displayEventFromChannel(fanInChannels(streamingCtx, eventChannels...))
}

func fanInChannels(commonCtx ctx.Context, channels ...<-chan cwl.StreamEvent) <-chan cwl.StreamEvent {
	var waitGroup sync.WaitGroup
	multiplexedChannel := make(chan cwl.StreamEvent)

	multiplexFunc := func(events <-chan cwl.StreamEvent) {
		defer waitGroup.Done()
		for event := range events {
			select {
			case <-commonCtx.Done():
				return
			case multiplexedChannel <- event:
			}
		}
	}

	waitGroup.Add(len(channels))
	for _, c := range channels {
		go multiplexFunc(c)
	}

	go func() {
		waitGroup.Wait()
		close(multiplexedChannel)
	}()

	return multiplexedChannel
}

func splitToBatchesBy(batchSize int, strs []string) [][]string {
	var batches [][]string
	totalStrings := len(strs)
	batchStart := 0
	for batchStart < totalStrings {
		batchEnd := batchStart + batchSize
		if batchEnd > totalStrings {
			batchEnd = totalStrings
		}
		batches = append(batches, strs[batchStart:batchEnd])
		batchStart = batchEnd
	}
	return batches
}

func (o *logsSharedOpts) displayLogGroup(logGroupName string, startTime, endTime *time.Time, filter string, streams ...string) error {
	output := o.cwlClient.GetLogsPaginated(cwl.GetLogsInput{
		LogGroupName: logGroupName,
		StartTime:    startTime,
		EndTime:      endTime,
		Filter:       filter,
		Streams:      streams,
	})
	for output.HasMoreLogs() {
		logs, err := output.NextLogs()
		if err != nil {
			return err
		}
		for _, line := range logs {
			printLn(line)
		}
	}
	return nil
}

func (o *logsSharedOpts) displayLogStreams(logGroupName string, startTime, endTime *time.Time, filter string, streams ...string) error {
	for _, stream := range streams {
		err := o.displayLogGroup(logGroupName, startTime, endTime, filter, stream)
		if err != nil {
			return err
		}
	}
	return nil
}
