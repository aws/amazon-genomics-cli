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
	LastOutput      string
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
	var keyToEventMap = make(map[string]ProgressEvent)
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(progressEvents))

	combinedStream := combineProgressEvents(progressEvents)

	for event := range combinedStream {
		if event.LastOutput != "" {
			log.Info().Msg(event.LastOutput)
		}
		keyToEventMap[event.UniqueKey] = event
	}

	return convertProgressEventsToResults(keyToEventMap, len(progressEvents))
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
	progressResult.UniqueKey = lastEvent.UniqueKey
}

func DisplayProgressBar(description string, progressEvents []ProgressStream) []Result {
	var keyToEventMap = make(map[string]ProgressEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	combinedStream := combineProgressEvents(progressEvents)

	barReceiver := runProgressBar(ctx, description, len(progressEvents))
	for event := range combinedStream {
		barReceiver <- event
		keyToEventMap[event.UniqueKey] = event
	}

	return convertProgressEventsToResults(keyToEventMap, len(progressEvents))
}

func convertProgressEventsToResults(keyToEventMap map[string]ProgressEvent, numberOfEvents int) []Result {
	var results = make([]Result, numberOfEvents)
	index := 0
	for _, progressResult := range keyToEventMap {
		results[index] = Result{
			progressResult.UniqueKey,
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
	for initialChannelOut := range channel {
		if initialChannelOut.Err != nil {
			stopProcessingEvent.UniqueKey = lastEvent.UniqueKey
			receiver <- stopProcessingEvent
			return
		} else {
			receiver <- initialChannelOut
			lastEvent = initialChannelOut
		}
	}
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
		var oldProgressEvents = make(map[string]ProgressEvent, numberOfChannels)
		var keyWithSteps = make(map[string]bool, numberOfChannels)
		totalSteps, currentStep := 0, 0
		for {
			select {
			case progressEvent := <-receiver:
				oldEvent, matchExists := oldProgressEvents[progressEvent.UniqueKey]

				if matchExists {
					totalSteps += progressEvent.TotalSteps - oldEvent.TotalSteps
					currentStep += progressEvent.CurrentStep - oldEvent.CurrentStep
					oldProgressEvents[oldEvent.UniqueKey] = progressEvent

					_, keysExist := keyWithSteps[progressEvent.UniqueKey]
					if !keysExist && progressEvent.TotalSteps > 0 {
						keyWithSteps[progressEvent.UniqueKey] = true
					}
				} else if progressEvent.UniqueKey != "" {
					totalSteps += progressEvent.TotalSteps
					currentStep += progressEvent.CurrentStep
					oldProgressEvents[progressEvent.UniqueKey] = progressEvent
				}

				if len(keyWithSteps) == numberOfChannels {
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
