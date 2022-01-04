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
	progressTemplate pb.ProgressBarTemplate = `{{ string . "description" }} [{{cycle . "o---" "-o--" "--o-" "---o" "--o-" "-o--" "o---" }}] {{ etime . }}`
)

type ProgressStream chan ProgressEvent

type ProgressEvent struct {
	CurrentStep     int
	TotalSteps      int
	StepDescription string
	Outputs         []string
	Err             error
	ExecutionName   string
	LastOutput      string
}

type Result struct {
	ExecutionName string
	Outputs       []string
	Err           error
}

func (p ProgressStream) DisplayProgress(description string) error {
	var lastEvent ProgressEvent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	numberOfChannels := 1
	barReceiver := runProgressBar(ctx, description, numberOfChannels)
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

func ShowExecution(progressStreams []ProgressStream) []Result {
	var keyToEventMap = make(map[string]ProgressEvent)

	combinedStream := combineProgressEvents(progressStreams)

	for event := range combinedStream {
		if event.LastOutput != "" {
			log.Info().Msg(event.LastOutput)
		}
		keyToEventMap[event.ExecutionName] = event
	}

	return convertProgressEventsToResults(keyToEventMap, len(progressStreams))
}

func updateResultFromStream(stream ProgressStream, progressResult *Result, wait *sync.WaitGroup) {
	defer wait.Done()
	var lastEvent ProgressEvent

	for event := range stream {
		if event.Err != nil {
			progressResult.Err = event.Err
		} else {
			lastEvent = event
		}
	}

	progressResult.Outputs = lastEvent.Outputs
	progressResult.ExecutionName = lastEvent.ExecutionName
}

func DisplayProgressBar(description string, progressEvents []ProgressStream) []Result {
	var keyToEventMap = make(map[string]ProgressEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	combinedStream := combineProgressEvents(progressEvents)

	barReceiver := runProgressBar(ctx, description, len(progressEvents))
	for event := range combinedStream {
		barReceiver <- event
		keyToEventMap[event.ExecutionName] = event
	}
	return convertProgressEventsToResults(keyToEventMap, len(progressEvents))
}

func convertProgressEventsToResults(keyToEventMap map[string]ProgressEvent, numberOfEvents int) []Result {
	var results = make([]Result, numberOfEvents)
	index := 0
	for _, progressResult := range keyToEventMap {
		results[index] = Result{
			progressResult.ExecutionName,
			progressResult.Outputs,
			progressResult.Err,
		}
		index++
	}

	return results
}

func combineProgressEvents(progressEventChannels []ProgressStream) ProgressStream {
	receiver := make(ProgressStream)
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(progressEventChannels))

	for _, channel := range progressEventChannels {
		go sendDataToReceiver(channel, &waitGroup, receiver)
	}

	go closeChannelAfterWaitGroup(receiver, &waitGroup)

	return receiver
}

func sendDataToReceiver(channel <-chan ProgressEvent, waitGroup *sync.WaitGroup, receiver chan ProgressEvent) {
	defer waitGroup.Done()

	var lastEvent ProgressEvent
	var stopProcessingEvent = ProgressEvent{
		CurrentStep: 1,
		TotalSteps:  1,
	}
	for cdkChannelOut := range channel {
		if cdkChannelOut.Err != nil {
			stopProcessingEvent.ExecutionName = lastEvent.ExecutionName
			stopProcessingEvent.Err = cdkChannelOut.Err
			stopProcessingEvent.Outputs = lastEvent.Outputs
			receiver <- stopProcessingEvent
			return
		} else {
			receiver <- cdkChannelOut
			lastEvent = cdkChannelOut
		}
	}
}

func closeChannelAfterWaitGroup(channel chan ProgressEvent, waitGroup *sync.WaitGroup) {
	waitGroup.Wait()
	close(channel)
}

func runProgressBar(ctx context.Context, description string, numberOfChannels int) chan ProgressEvent {
	receiver := make(chan ProgressEvent)
	bar := progressTemplate.New(0).
		Set("description", description).
		Start()

	go func() {
		var oldProgressEvents = make(map[string]ProgressEvent, numberOfChannels)
		var keyWithSteps = make(map[string]bool, numberOfChannels)
		totalSteps, currentStep := 0, 0
		for {
			select {
			case progressEvent := <-receiver:
				oldEvent, matchExists := oldProgressEvents[progressEvent.ExecutionName]

				if matchExists {
					totalSteps += progressEvent.TotalSteps - oldEvent.TotalSteps
					currentStep += progressEvent.CurrentStep - oldEvent.CurrentStep
					oldProgressEvents[oldEvent.ExecutionName] = progressEvent

					_, keysExist := keyWithSteps[progressEvent.ExecutionName]
					if !keysExist && progressEvent.TotalSteps > 0 {
						keyWithSteps[progressEvent.ExecutionName] = true
					}
				} else if progressEvent.ExecutionName != "" {
					totalSteps += progressEvent.TotalSteps
					currentStep += progressEvent.CurrentStep
					oldProgressEvents[progressEvent.ExecutionName] = progressEvent
				}

				if len(keyWithSteps) == numberOfChannels {
					bar.SetCurrent(int64(currentStep))
				}

				if progressEvent.StepDescription != "" {
					bar.Set("description", progressEvent.StepDescription)
				} else {
					bar.Set("description", description)
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
