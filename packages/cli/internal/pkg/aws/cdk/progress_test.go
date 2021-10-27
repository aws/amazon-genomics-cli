package cdk

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SendDataToReceiverAndUpdateResult_Success(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testChannel, receivingChannel := make(ProgressStream), make(ProgressStream)
	progressResult := Result{}

	go sendDataToReceiverAndUpdateResult(testChannel, &progressResult, &waitGroup, receivingChannel)

	sentEvent := ProgressEvent{UniqueKey: "myKey", Outputs: []string{"hi"}}
	testChannel <- sentEvent
	close(testChannel)

	channelOutput := <-receivingChannel

	waitGroup.Wait()

	expectedProgressResult := Result{
		UniqueKey: "myKey",
		Outputs:   []string{"hi"},
	}

	assert.Equal(t, expectedProgressResult, progressResult)
	assert.Equal(t, sentEvent, channelOutput)
}

func Test_SendDataToReceiverAndUpdateResult_Error(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testChannel, receivingChannel := make(ProgressStream), make(ProgressStream)
	progressResult := Result{}

	go sendDataToReceiverAndUpdateResult(testChannel, &progressResult, &waitGroup, receivingChannel)

	sentEvent := ProgressEvent{UniqueKey: "someKey", Outputs: []string{"hi"}}
	sentErrorEvent := ProgressEvent{Err: errors.New("some error"), UniqueKey: "someKey"}
	testChannel <- sentEvent
	channelOutput := <-receivingChannel
	testChannel <- sentErrorEvent
	channelOutput = <-receivingChannel
	close(testChannel)

	waitGroup.Wait()

	expectedProgressResult := Result{
		UniqueKey: "someKey",
		Outputs:   []string{"hi"},
		Err:       errors.New("some error"),
	}

	assert.Equal(t, expectedProgressResult, progressResult)

	expectedEvent := ProgressEvent{UniqueKey: "someKey", CurrentStep: 0, TotalSteps: 0}
	assert.Equal(t, expectedEvent, channelOutput)
}

func Test_updateResultFromStream_Success(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testChannel := make(ProgressStream)
	progressResult := Result{}

	go updateResultFromStream(testChannel, &progressResult, &waitGroup)

	sentEvent := ProgressEvent{UniqueKey: "someKey", Outputs: []string{"hi"}}
	testChannel <- sentEvent
	close(testChannel)

	waitGroup.Wait()

	expectedProgressResult := Result{
		UniqueKey: "someKey",
		Outputs:   []string{"hi"},
	}

	assert.Equal(t, expectedProgressResult, progressResult)
}

func Test_updateResultFromStream_Error(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testChannel := make(ProgressStream)
	progressResult := Result{}

	go updateResultFromStream(testChannel, &progressResult, &waitGroup)

	sentEvent := ProgressEvent{UniqueKey: "someKey", Outputs: []string{"hi"}}
	testChannel <- sentEvent

	sentErrorEvent := ProgressEvent{
		Err:       errors.New("some error"),
		UniqueKey: "someKey",
		Outputs:   []string{"hi"},
	}
	testChannel <- sentErrorEvent
	close(testChannel)

	waitGroup.Wait()

	expectedProgressResult := Result{
		UniqueKey: "someKey",
		Outputs:   []string{"hi"},
		Err:       errors.New("some error"),
	}

	assert.Equal(t, expectedProgressResult, progressResult)
}
