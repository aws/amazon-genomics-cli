package cli

import (
	ctx "context"
	"fmt"
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
	logFilterFlagDescription = "Match terms, phrases, or values in the logs."

	tailFlag            = "tail"
	tailFlagShort       = "t"
	tailFlagDescription = "Follow the log output."
)

var printLn = fmt.Println

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

func (o *logsSharedOpts) followLogGroup(logGroupName string, streams ...string) error {
	stream := o.cwlClient.StreamLogs(ctx.Background(), logGroupName, streams...)
	for event := range stream {
		if event.Err != nil {
			return event.Err
		}
		if len(event.Logs) > 0 {
			for _, line := range event.Logs {
				_, _ = printLn(line)
			}
		} else {
			log.Debug().Msg("No new logs")
		}
	}
	return nil
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
			_, _ = printLn(line)
		}
	}
	return nil
}
