package cli

import (
	ctx "context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	testLogGroupName = "test-log-group-name"
)

func TestLogsEngineOpts_Validate_LookBackFlag(t *testing.T) {
	opts := logsEngineOpts{logsEngineVars: logsEngineVars{logsSharedVars{lookBack: "5h"}}}
	before := time.Now().Add(-5 * time.Hour)
	err := opts.Validate()
	after := time.Now().Add(-5 * time.Hour)
	assert.NoError(t, err)
	assert.True(t, before.Before(*opts.startTime))
	assert.True(t, after.After(*opts.startTime))
}

func TestLogsEngineOpts_Validate_LookBackError(t *testing.T) {
	opts := logsEngineOpts{logsEngineVars: logsEngineVars{logsSharedVars{lookBack: "abc"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("time: invalid duration \"abc\""), err)
}

func TestLogsEngineOpts_Validate_StartEndFlags(t *testing.T) {
	start := time.Unix(0, 773391600000)
	end := time.Unix(0, 773391700000)
	opts := logsEngineOpts{logsEngineVars: logsEngineVars{logsSharedVars{startString: start.Format(time.RFC3339Nano), endString: end.Format(time.RFC3339Nano)}}}
	err := opts.Validate()
	assert.NoError(t, err)
	assert.True(t, start.Equal(*opts.startTime))
	assert.True(t, end.Equal(*opts.endTime))
}

func TestLogsEngineOpts_Validate_StartError(t *testing.T) {
	opts := logsEngineOpts{logsEngineVars: logsEngineVars{logsSharedVars{startString: "abc"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("Could not find format for \"abc\""), err)
}

func TestLogsEngineOpts_Validate_EndError(t *testing.T) {
	opts := logsEngineOpts{logsEngineVars: logsEngineVars{logsSharedVars{endString: "abc"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("Could not find format for \"abc\""), err)
}

func TestLogsEngineOpts_Validate_FlagConflictError(t *testing.T) {
	opts := logsEngineOpts{logsEngineVars: logsEngineVars{logsSharedVars{startString: "1/1/1990", lookBack: "1h"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("a look back period cannot be specified together with start or end times"), err)
}

func TestLogsEngineOpts_Execute_Group(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	logPaginatorMock := awsmocks.NewMockCwlLogPaginator(ctrl)
	opts := logsEngineOpts{
		logsSharedOpts: logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsEngineVars: logsEngineVars{logsSharedVars{contextName: testContextName1}},
	}

	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{EngineLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().GetLogsPaginated(cwl.GetLogsInput{testLogGroupName, nil, nil, "", nil}).Return(logPaginatorMock)
	gomock.InOrder(logPaginatorMock.EXPECT().HasMoreLogs().Return(true), logPaginatorMock.EXPECT().HasMoreLogs().Return(false))
	logPaginatorMock.EXPECT().NextLogs().Return([]string{"log"}, nil)

	err := opts.Execute()
	assert.NoError(t, err)
}

func TestLogsEngineOpts_Execute_InfoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsEngineOpts{
		logsSharedOpts: logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsEngineVars: logsEngineVars{logsSharedVars{contextName: testContextName1}},
	}

	someErr := fmt.Errorf("some info error")
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{}, someErr)

	err := opts.Execute()
	assert.Equal(t, someErr, err)
}

func TestLogsEngineOpts_Execute_LogError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	logPaginatorMock := awsmocks.NewMockCwlLogPaginator(ctrl)
	opts := logsEngineOpts{
		logsSharedOpts: logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsEngineVars: logsEngineVars{logsSharedVars{contextName: testContextName1}},
	}

	someErr := fmt.Errorf("some log error")
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{EngineLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().GetLogsPaginated(cwl.GetLogsInput{testLogGroupName, nil, nil, "", nil}).Return(logPaginatorMock)
	logPaginatorMock.EXPECT().HasMoreLogs().Return(true)
	logPaginatorMock.EXPECT().NextLogs().Return(nil, someErr)

	err := opts.Execute()
	assert.Equal(t, someErr, err)
}

func TestLogsEngineOpts_Execute_Stream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsEngineOpts{
		logsSharedOpts: logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsEngineVars: logsEngineVars{logsSharedVars{contextName: testContextName1, tail: true}},
	}
	stream := make(chan cwl.StreamEvent)
	go func() { stream <- cwl.StreamEvent{Logs: []string{"log"}}; close(stream) }()
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{EngineLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().StreamLogs(ctx.Background(), testLogGroupName).Return(stream)

	err := opts.Execute()
	assert.NoError(t, err)
}

func TestLogsEngineOpts_Execute_StreamError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsEngineOpts{
		logsSharedOpts: logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsEngineVars: logsEngineVars{logsSharedVars{contextName: testContextName1, tail: true}},
	}
	someErr := fmt.Errorf("some stream error")
	stream := make(chan cwl.StreamEvent)
	go func() { stream <- cwl.StreamEvent{Err: someErr} }()
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{EngineLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().StreamLogs(ctx.Background(), testLogGroupName).Return(stream)

	err := opts.Execute()
	assert.Equal(t, someErr, err)
}
