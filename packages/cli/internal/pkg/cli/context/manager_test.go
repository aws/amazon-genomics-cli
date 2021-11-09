package context

import (
	"errors"
	"testing"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestManager_ProcessExecution(t *testing.T) {
	origVerbose := logging.Verbose
	origShowExecution := showExecution
	defer func() {
		logging.Verbose = origVerbose
		showExecution = origShowExecution
	}()
	logging.Verbose = true

	mockClients := createMocks(t)

	showExecution = mockClients.cdkMock.ShowExecution
	cdkResults := []cdk.Result{{Outputs: []string{"some message"}, UniqueKey: testContextName1}, {UniqueKey: testContextName2, Err: errors.New("Some error")}}
	progressStreams := []cdk.ProgressStream{mockClients.progressStream1, mockClients.progressStream2}
	mockClients.cdkMock.EXPECT().ShowExecution(progressStreams).Return(cdkResults)
	defer close(mockClients.progressStream1)
	defer close(mockClients.progressStream2)

	defer mockClients.ctrl.Finish()
	manager := Manager{
		Cdk:       mockClients.cdkMock,
		Project:   mockClients.projMock,
		Ssm:       mockClients.ssmMock,
		Config:    mockClients.configMock,
		Cfn:       mockClients.cfnMock,
		baseProps: baseProps{homeDir: testHomeDir},
	}

	manager.processExecution(progressStreams, "some description")

	expectedProgressResults := []ProgressResult{
		{
			Context: testContextName1,
			Outputs: []string{"some message"},
		},
		{
			Context: testContextName2,
			Err:     errors.New("Some error"),
		},
	}
	assert.ElementsMatch(t, expectedProgressResults, manager.progressResults)
}
