package cdk

import (
	"context"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/rs/zerolog/log"
)

var (
	sleep                                   = time.Sleep
	progressTemplate pb.ProgressBarTemplate = `{{ string . "description" }} {{ bar . }}{{ etime . }}`
)

type ProgressStream chan ProgressEvent

type ProgressEvent struct {
	CurrentStep     int
	TotalSteps      int
	StepDescription string
	Outputs         []string
	Err             error
}

func (p ProgressStream) DisplayProgress(description string) error {
	var lastEvent ProgressEvent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	barReceiver := runProgressBar(ctx, description)
	for event := range p {
		barReceiver <- event
		if event.Err != nil {
			for _, line := range lastEvent.Outputs {
				log.Error().Msg(line)
			}
			return event.Err
		}
		lastEvent = event
	}
	return nil
}

func runProgressBar(ctx context.Context, description string) chan ProgressEvent {
	receiver := make(chan ProgressEvent)
	bar := progressTemplate.New(1).
		Set("description", description).
		Start()

	go func() {
		for {
			select {
			case event := <-receiver:
				if event.TotalSteps > 0 {
					bar.SetTotal(int64(event.TotalSteps))
					bar.SetCurrent(int64(event.CurrentStep))
				}
			case <-ctx.Done():
				close(receiver)
				bar.SetCurrent(bar.Total()).Finish()
				return
			default:
				// No new events; nothing to do
			}
			sleep(100 * time.Millisecond)
		}
	}()

	return receiver
}
