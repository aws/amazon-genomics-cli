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

	go sendDataToReceiver(testChannel, &waitGroup, receivingChannel)

	sentEvent := ProgressEvent{ExecutionName: "myKey", Outputs: []string{"hi"}}
	testChannel <- sentEvent
	close(testChannel)

	channelOutput := <-receivingChannel

	waitGroup.Wait()

	assert.Equal(t, sentEvent, channelOutput)
}

func Test_SendDataToReceiverAndUpdateResult_Error(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testChannel, receivingChannel := make(ProgressStream), make(ProgressStream)

	go sendDataToReceiver(testChannel, &waitGroup, receivingChannel)

	sentEvent := ProgressEvent{ExecutionName: "someKey", Outputs: []string{"hi"}}
	sentErrorEvent := ProgressEvent{Err: errors.New("some error"), ExecutionName: "someKey"}
	testChannel <- sentEvent
	<-receivingChannel
	testChannel <- sentErrorEvent
	channelOutput := <-receivingChannel
	close(testChannel)

	waitGroup.Wait()

	expectedEvent := ProgressEvent{ExecutionName: "someKey", CurrentStep: 1, TotalSteps: 1}
	assert.Equal(t, expectedEvent, channelOutput)
}

func Test_updateResultFromStream_Success(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testChannel := make(ProgressStream)
	progressResult := Result{}

	go updateResultFromStream(testChannel, &progressResult, &waitGroup)

	sentEvent := ProgressEvent{ExecutionName: "someKey", Outputs: []string{"hi"}}
	testChannel <- sentEvent
	close(testChannel)

	waitGroup.Wait()

	expectedProgressResult := Result{
		ExecutionName: "someKey",
		Outputs:       []string{"hi"},
	}

	assert.Equal(t, expectedProgressResult, progressResult)
}

func Test_updateResultFromStream_Error(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	testChannel := make(ProgressStream)
	progressResult := Result{}

	go updateResultFromStream(testChannel, &progressResult, &waitGroup)

	sentEvent := ProgressEvent{ExecutionName: "someKey", Outputs: []string{"hi"}}
	testChannel <- sentEvent

	sentErrorEvent := ProgressEvent{
		Err:           errors.New("some error"),
		ExecutionName: "someKey",
		Outputs:       []string{"hi"},
	}
	testChannel <- sentErrorEvent
	close(testChannel)

	waitGroup.Wait()

	expectedProgressResult := Result{
		ExecutionName: "someKey",
		Outputs:       []string{"hi"},
		Err:           errors.New("some error"),
	}

	assert.Equal(t, expectedProgressResult, progressResult)
}
