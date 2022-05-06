package cli

import (
	ctx "context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/context"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogsAdapterOpts_Validate_LookBackFlag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1, lookBack: "5h"}},
	}
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	before := time.Now().Add(-5 * time.Hour)
	err := opts.Validate()
	after := time.Now().Add(-5 * time.Hour)
	assert.NoError(t, err)
	assert.True(t, before.Before(*opts.startTime))
	assert.True(t, after.After(*opts.startTime))
}

func TestLogsAdapterOpts_Validate_LookBackError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1, lookBack: "abc"}},
	}
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("time: invalid duration \"abc\""), err)
}

func TestLogsAdapterOpts_Validate_StartEndFlags(t *testing.T) {
	start := time.Unix(0, 773391600000)
	end := time.Unix(0, 773391700000)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts: logsSharedOpts{ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{
			logsSharedVars{
				contextName: testContextName1, startString: start.Format(time.RFC3339Nano), endString: end.Format(time.RFC3339Nano),
			},
		},
	}
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	err := opts.Validate()
	assert.NoError(t, err)
	assert.True(t, start.Equal(*opts.startTime))
	assert.True(t, end.Equal(*opts.endTime))
}

func TestLogsAdapterOpts_Validate_StartError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1, startString: "abc"}},
	}
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("Could not find format for \"abc\""), err)
}

func TestLogsAdapterOpts_Validate_EndError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1, endString: "abc"}},
	}
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("Could not find format for \"abc\""), err)
}

func TestLogsAdapterOpts_Validate_FlagConflictError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1, startString: "1/1/1990", lookBack: "1h"}},
	}
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("a look back period cannot be specified together with start or end times"), err)
}

func TestLogsAdapterOpts_Validate_ToilError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars: logsSharedVars{contextName: testContextName1}},
	}
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.TOIL}}}}, nil).AnyTimes()

	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("Contexts using the toil engine do not have adapters to collect logs from"), err)
}

func TestLogsAdapterOpts_Execute_Group(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	logPaginatorMock := awsmocks.NewMockCwlLogPaginator(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1}},
	}

	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{WesLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().GetLogsPaginated(cwl.GetLogsInput{LogGroupName: testLogGroupName}).Return(logPaginatorMock)
	gomock.InOrder(logPaginatorMock.EXPECT().HasMoreLogs().Return(true), logPaginatorMock.EXPECT().HasMoreLogs().Return(false))
	logPaginatorMock.EXPECT().NextLogs().Return([]string{"log"}, nil)

	err := opts.Execute()
	assert.NoError(t, err)
}

func TestLogsAdapterOpts_Execute_InfoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1}},
	}

	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	someErr := fmt.Errorf("some info error")
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{}, someErr)

	err := opts.Execute()
	assert.Equal(t, someErr, err)
}

func TestLogsAdapterOpts_Execute_LogError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	logPaginatorMock := awsmocks.NewMockCwlLogPaginator(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1}},
	}

	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	someErr := fmt.Errorf("some log error")
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{WesLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().GetLogsPaginated(cwl.GetLogsInput{LogGroupName: testLogGroupName}).Return(logPaginatorMock)
	logPaginatorMock.EXPECT().HasMoreLogs().Return(true)
	logPaginatorMock.EXPECT().NextLogs().Return(nil, someErr)

	err := opts.Execute()
	assert.Equal(t, someErr, err)
}

func TestLogsAdapterOpts_Execute_Stream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1, tail: true}},
	}
	stream := make(chan cwl.StreamEvent)
	go func() { stream <- cwl.StreamEvent{Logs: []string{"log"}}; close(stream) }()
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{WesLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().StreamLogs(ctx.Background(), testLogGroupName).Return(stream)

	err := opts.Execute()
	assert.NoError(t, err)
}

func TestLogsAdapterOpts_Execute_StreamError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cwlMock := awsmocks.NewMockCwlClient(ctrl)
	ctxMock := contextmocks.NewMockContextManager(ctrl)
	opts := logsAdapterOpts{
		logsSharedOpts:  logsSharedOpts{cwlClient: cwlMock, ctxManager: ctxMock},
		logsAdapterVars: logsAdapterVars{logsSharedVars{contextName: testContextName1, tail: true}},
	}
	someErr := fmt.Errorf("some stream error")
	stream := make(chan cwl.StreamEvent)
	go func() { stream <- cwl.StreamEvent{Err: someErr} }()
	ctxMock.EXPECT().List().Return(map[string]context.Summary{testContextName1: {Engines: []spec.Engine{{Engine: constants.CROMWELL}}}}, nil).AnyTimes()
	ctxMock.EXPECT().Info(testContextName1).Return(context.Detail{WesLogGroupName: testLogGroupName}, nil)
	cwlMock.EXPECT().StreamLogs(ctx.Background(), testLogGroupName).Return(stream)

	err := opts.Execute()
	assert.Equal(t, someErr, err)
}
