package cli

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/batch"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/workflow"
	awsmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/aws"
	contextmocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/context"
	managermocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogsWorkflowOpts_Validate_LookBackFlag(t *testing.T) {
	opts := logsWorkflowOpts{logsWorkflowVars: logsWorkflowVars{logsSharedVars: logsSharedVars{lookBack: "5h"}}}
	before := time.Now().Add(-5 * time.Hour)
	err := opts.Validate()
	after := time.Now().Add(-5 * time.Hour)
	assert.NoError(t, err)
	assert.True(t, before.Before(*opts.startTime))
	assert.True(t, after.After(*opts.startTime))
}

func TestLogsWorkflowOpts_Validate_LookBackError(t *testing.T) {
	opts := logsWorkflowOpts{logsWorkflowVars: logsWorkflowVars{logsSharedVars: logsSharedVars{lookBack: "abc"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("time: invalid duration \"abc\""), err)
}

func TestLogsWorkflowOpts_Validate_StartEndFlags(t *testing.T) {
	start := time.Unix(0, 773391600000)
	end := time.Unix(0, 773391700000)
	opts := logsWorkflowOpts{logsWorkflowVars: logsWorkflowVars{logsSharedVars: logsSharedVars{startString: start.Format(time.RFC3339Nano), endString: end.Format(time.RFC3339Nano)}}}
	err := opts.Validate()
	assert.NoError(t, err)
	assert.True(t, start.Equal(*opts.startTime))
	assert.True(t, end.Equal(*opts.endTime))
}

func TestLogsWorkflowOpts_Validate_StartError(t *testing.T) {
	opts := logsWorkflowOpts{logsWorkflowVars: logsWorkflowVars{logsSharedVars: logsSharedVars{startString: "abc"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("Could not find format for \"abc\""), err)
}

func TestLogsWorkflowOpts_Validate_EndError(t *testing.T) {
	opts := logsWorkflowOpts{logsWorkflowVars: logsWorkflowVars{logsSharedVars: logsSharedVars{endString: "abc"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("Could not find format for \"abc\""), err)
}

func TestLogsWorkflowOpts_Validate_FlagConflictError(t *testing.T) {
	opts := logsWorkflowOpts{logsWorkflowVars: logsWorkflowVars{logsSharedVars: logsSharedVars{startString: "1/1/1990", lookBack: "1h"}}}
	err := opts.Validate()
	assert.Equal(t, fmt.Errorf("a look back period cannot be specified together with start or end times"), err)
}

func Test_filterCachedJobIds(t *testing.T) {
	tests := map[string]struct {
		ids      []string
		expected []string
	}{
		"empty": {
			ids:      []string{},
			expected: []string(nil),
		},
		"allNotCached": {
			ids:      []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		"allCached": {
			ids:      []string{cachedJobId, cachedJobId, cachedJobId},
			expected: []string(nil),
		},
		"mix": {
			ids:      []string{"a", cachedJobId, cachedJobId, "b", "c", cachedJobId},
			expected: []string{"a", "b", "c"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := filterCachedJobIds(tt.ids)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLogsWorkflowOpts_Execute(t *testing.T) {

	const (
		testLogGroupName  = "/aws/batch/job"
		testWorkflowName  = "Test Workflow"
		testRunId         = "Test Workflow Run Id"
		testTaskName      = "Test Task Name"
		testJobId         = "Test Job Id"
		testJobName       = "Test Job Name"
		testLogStreamName = "Test Log Stream Name"
		testLogPage1      = "Test Log Page 1"
	)

	testTask := workflow.Task{
		Name:  testTaskName,
		JobId: testJobId,
	}

	testRunLog := workflow.RunLog{
		RunId: testRunId,
		State: "COMPLETE",
		Tasks: []workflow.Task{testTask},
	}

	testCachedTask := workflow.Task{
		Name:  testTaskName,
		JobId: cachedJobId,
	}

	testJob := batch.Job{
		JobId:         testJobId,
		JobName:       testJobName,
		LogStreamName: testLogStreamName,
	}

	testInstanceSummary := workflow.InstanceSummary{
		Id:           testRunId,
		WorkflowName: testWorkflowName,
		ContextName:  testContextName,
	}

	tests := map[string]struct {
		setupOps             func(*logsWorkflowOpts, *awsmocks.MockCwlLogPaginator)
		expectedOutput       string
		expectedErrorMessage string
	}{
		"runId empty log": {
			setupOps: func(opts *logsWorkflowOpts, cwlLopPaginator *awsmocks.MockCwlLogPaginator) {
				opts.workflowName = testWorkflowName
				opts.runId = testRunId
				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					GetRunLog(testRunId).Return(testRunLog, nil)
			},
			expectedOutput: "RunId: Test Workflow Run Id\nState: COMPLETE\nTasks: \n\tName\t\tJobId\t\tStartTime\tStopTimeExitCode\n\tTest Task Name\tTest Job Id\t<nil>\t\t<nil>\t\n\t\n",
		},
		"runId no jobs": {
			setupOps: func(opts *logsWorkflowOpts, cwlLopPaginator *awsmocks.MockCwlLogPaginator) {
				opts.workflowName = testWorkflowName
				opts.runId = testRunId
				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					GetRunLog(testRunId).Return(workflow.RunLog{
					RunId: testRunId,
					State: "COMPLETE",
					Tasks: []workflow.Task(nil),
				}, nil)
			},
			expectedOutput: "RunId: Test Workflow Run Id\nState: COMPLETE\nTasks: No task logs available\n",
		},
		"runId empty cached": {
			setupOps: func(opts *logsWorkflowOpts, cwlLopPaginator *awsmocks.MockCwlLogPaginator) {
				opts.workflowName = testWorkflowName
				opts.runId = testRunId

				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					GetRunLog(testRunId).Return(workflow.RunLog{
					RunId: testRunId,
					State: "COMPLETE",
					Tasks: []workflow.Task{testCachedTask},
				}, nil)
			},
			expectedOutput: "RunId: Test Workflow Run Id\nState: COMPLETE\nTasks: \n\tName\t\tJobId\tStartTime\tStopTimeExitCode\n\tTest Task Name\tXXXXX\t<nil>\t\t<nil>\t\n\t\n",
		},
		"workflow name single task empty log": {
			setupOps: func(opts *logsWorkflowOpts, cwlLopPaginator *awsmocks.MockCwlLogPaginator) {
				opts.workflowName = testWorkflowName
				opts.taskId = testJobId

				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					GetRunLog(testRunId).Return(workflow.RunLog{
					RunId: testRunId,
					State: "COMPLETE",
					Tasks: []workflow.Task{testTask},
				}, nil)
				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					StatusWorkflowByName(testWorkflowName, 1).Return([]workflow.InstanceSummary{testInstanceSummary}, nil)
				opts.batchClient.(*awsmocks.MockBatchClient).EXPECT().
					GetJobs([]string{testJobId}).Return([]batch.Job{testJob}, nil)
				opts.cwlClient.(*awsmocks.MockCwlClient).EXPECT().
					GetLogsPaginated(cwl.GetLogsInput{
						LogGroupName: testLogGroupName,
						StartTime:    nil,
						EndTime:      nil,
						Filter:       "",
						Streams:      []string{testLogStreamName},
					}).Return(cwlLopPaginator)
				cwlLopPaginator.EXPECT().HasMoreLogs().Return(false)
			},
			expectedOutput: "",
		},
		"workflow name no runs": {
			setupOps: func(opts *logsWorkflowOpts, cwlLopPaginator *awsmocks.MockCwlLogPaginator) {
				opts.workflowName = testWorkflowName

				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					StatusWorkflowByName(testWorkflowName, 1).Return([]workflow.InstanceSummary(nil), nil)
			},
			expectedOutput: "",
		},
		"runId one page": {
			setupOps: func(opts *logsWorkflowOpts, cwlLopPaginator *awsmocks.MockCwlLogPaginator) {
				opts.workflowName = testWorkflowName
				opts.runId = testRunId
				opts.allTasks = true

				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					GetRunLog(testRunId).Return(workflow.RunLog{
					RunId: testRunId,
					State: "COMPLETE",
					Tasks: []workflow.Task{testTask},
				}, nil)
				opts.batchClient.(*awsmocks.MockBatchClient).EXPECT().
					GetJobs([]string{testJobId}).Return([]batch.Job{testJob}, nil)
				opts.cwlClient.(*awsmocks.MockCwlClient).EXPECT().
					GetLogsPaginated(cwl.GetLogsInput{
						LogGroupName: testLogGroupName,
						StartTime:    nil,
						EndTime:      nil,
						Filter:       "",
						Streams:      []string{testLogStreamName},
					}).Return(cwlLopPaginator)
				cwlLopPaginator.EXPECT().HasMoreLogs().Return(true)
				cwlLopPaginator.EXPECT().HasMoreLogs().Return(false)
				cwlLopPaginator.EXPECT().NextLogs().Return([]string{testLogPage1}, nil)
			},
			expectedOutput: "Test Log Page 1\n",
		},
		"receives error": {
			setupOps: func(opts *logsWorkflowOpts, cwlLopPaginator *awsmocks.MockCwlLogPaginator) {
				opts.workflowName = testWorkflowName
				opts.runId = testRunId
				opts.allTasks = true

				opts.workflowManager.(*managermocks.MockWorkflowManager).EXPECT().
					GetRunLog(testRunId).Return(workflow.RunLog{
					RunId: testRunId,
					State: "COMPLETE",
					Tasks: []workflow.Task{testTask},
				}, nil)
				opts.batchClient.(*awsmocks.MockBatchClient).EXPECT().
					GetJobs([]string{testJobId}).Return([]batch.Job{testJob}, nil)
				opts.cwlClient.(*awsmocks.MockCwlClient).EXPECT().
					GetLogsPaginated(cwl.GetLogsInput{
						LogGroupName: testLogGroupName,
						StartTime:    nil,
						EndTime:      nil,
						Filter:       "",
						Streams:      []string{testLogStreamName},
					}).Return(cwlLopPaginator)
				cwlLopPaginator.EXPECT().HasMoreLogs().Return(true)
				cwlLopPaginator.EXPECT().NextLogs().Return([]string{}, errors.New("some error"))
			},
			expectedOutput:       "",
			expectedErrorMessage: "some error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var output strings.Builder
			origPrintLn := printLn
			printLn = func(args ...interface{}) {
				output.WriteString(fmt.Sprintln(args...))
			}
			defer func() { printLn = origPrintLn }()

			mockWorkflow := managermocks.NewMockWorkflowManager(ctrl)
			mockContext := contextmocks.NewMockContextManager(ctrl)
			mockCwl := awsmocks.NewMockCwlClient(ctrl)
			mockBatch := awsmocks.NewMockBatchClient(ctrl)
			vars := logsWorkflowVars{logsSharedVars: logsSharedVars{}}
			opts := &logsWorkflowOpts{
				logsWorkflowVars: vars,
				logsSharedOpts: logsSharedOpts{
					ctxManager: mockContext,
					cwlClient:  mockCwl,
				},
				batchClient:     mockBatch,
				workflowManager: mockWorkflow,
			}
			cwlPager := awsmocks.NewMockCwlLogPaginator(ctrl)
			tt.setupOps(opts, cwlPager)
			err := opts.Execute()

			if tt.expectedErrorMessage != "" {
				assert.Error(t, err, tt.expectedErrorMessage)
			} else {
				assert.NoError(t, err)
				actualOutput := output.String()
				assert.Equal(t, tt.expectedOutput, actualOutput)
			}
		})
	}
}
