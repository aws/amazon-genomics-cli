package cdk

import (
	"context"
	"sync"
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
	UniqueKey       string
}

type Result struct {
	UniqueKey string
	Outputs   []string
	Err       error
}

func (p ProgressStream) DisplayProgress(description string) error {
	var lastEvent ProgressEvent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	barReceiver := runProgressBar(ctx, description, 1)
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

func ShowExecution(progressEvents []ProgressStream) []Result {
	var progressResults []*Result
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(progressEvents))

	for _, progressStream := range progressEvents {
		progressResult := &Result{}
		progressResults = append(progressResults, progressResult)
		go updateResultFromStream(progressStream, progressResult, &waitGroup)
	}

	waitGroup.Wait()

	var results []Result
	for _, progressResult := range progressResults {
		results = append(results, *progressResult)
	}

	return results
}

func updateResultFromStream(stream ProgressStream, progressResult *Result, wait *sync.WaitGroup) {
	var lastEvent ProgressEvent

	for event := range stream {
		if event.Err != nil {
			progressResult.Err = event.Err
		} else {
			lastEvent = event
		}
	}

	progressResult.Outputs = lastEvent.Outputs
	progressResult.UniqueKey = lastEvent.UniqueKey
	wait.Done()
}

func (client Client) DisplayProgressBar(description string, progressEvents []ProgressStream) []Result {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	singleStream, progressResults := combineProgressEvents(progressEvents)

	barReceiver := runProgressBar(ctx, description, len(progressEvents))
	for event := range singleStream {
		barReceiver <- event
	}

	var results []Result
	for _, progressResult := range progressResults {
		results = append(results, *progressResult)
	}

	return results
}

func combineProgressEvents(progressEventChannels []ProgressStream) (ProgressStream, []*Result) {
	var results []*Result
	receiver := make(ProgressStream)
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(progressEventChannels))

	for _, channel := range progressEventChannels {
		progressResult := &Result{}
		results = append(results, progressResult)
		go sendDataToReceiverAndUpdateResult(channel, progressResult, &waitGroup, receiver)
	}

	go closeChannelAfterWaitGroup(receiver, &waitGroup)

	return receiver, results
}

func sendDataToReceiverAndUpdateResult(channel <-chan ProgressEvent, progressResult *Result, waitGroup *sync.WaitGroup, receiver chan ProgressEvent) {
	defer waitGroup.Done()

	var lastEvent ProgressEvent
	var stopProcessingEvent = ProgressEvent{
		CurrentStep: 0,
		TotalSteps:  0,
	}
	for initialChannelOut := range channel {
		if initialChannelOut.Err != nil {
			progressResult.Err = initialChannelOut.Err

			stopProcessingEvent.UniqueKey = lastEvent.UniqueKey
			receiver <- stopProcessingEvent

			progressResult.Outputs = lastEvent.Outputs
			progressResult.UniqueKey = lastEvent.UniqueKey
			return
		} else {
			receiver <- initialChannelOut
			lastEvent = initialChannelOut
		}
	}

	progressResult.Outputs = lastEvent.Outputs
	progressResult.UniqueKey = lastEvent.UniqueKey
}

func closeChannelAfterWaitGroup(channel chan ProgressEvent, waitGroup *sync.WaitGroup) {
	waitGroup.Wait()
	close(channel)
}

func runProgressBar(ctx context.Context, description string, numberOfChannels int) chan ProgressEvent {
	receiver := make(chan ProgressEvent)
	bar := progressTemplate.New(1).
		Set("description", description).
		Start()

	go func() {
		var oldProgressEvents = make(map[string]ProgressEvent, numberOfChannels*2)
		totalSteps, currentStep := 0, 0
		for {
			select {
			case progressEvent := <-receiver:
				oldEvent, matchExists := oldProgressEvents[progressEvent.UniqueKey]

				if matchExists {
					totalSteps += progressEvent.TotalSteps - oldEvent.TotalSteps
					currentStep += progressEvent.CurrentStep - oldEvent.CurrentStep
					oldProgressEvents[oldEvent.UniqueKey] = progressEvent
				} else if progressEvent.UniqueKey != "" {
					totalSteps += progressEvent.TotalSteps
					currentStep += progressEvent.CurrentStep
					oldProgressEvents[progressEvent.UniqueKey] = progressEvent
				}

				if totalSteps > 0 {
					bar.SetTotal(int64(totalSteps))
					bar.SetCurrent(int64(currentStep))
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
