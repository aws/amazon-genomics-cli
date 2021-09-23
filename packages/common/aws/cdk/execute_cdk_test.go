package cdk

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	iomocks "github.com/aws/amazon-genomics-cli/common/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testExecuteCommandSuccessArg   = "test-execute-command-success-arg"
	testExecuteCommandFailureArg   = "test-execute-command-failure-arg"
	testExecuteCommandProgressLine = "  3/10 |4:56:17 PM | CREATE_COMPLETE      | AWS::IAM::Policy               | TaskBatch/BatchRole/DefaultPolicy (TaskBatchBatchRoleDefaultPolicyB9AAE3A1)"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestExecuteCdkCommand_Success(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osRemoveAll = mockOs.RemoveAll
	mockOs.EXPECT().RemoveAll(gomock.Any()).Return(nil).Times(0)

	progressStream, _ := executeCdkCommand(t.TempDir(), []string{testExecuteCommandSuccessArg})
	event1 := <-progressStream
	assert.Equal(t, 3, event1.CurrentStep)
	assert.Equal(t, 10, event1.TotalSteps)
	assert.Equal(t, testExecuteCommandProgressLine, event1.Outputs[0])
	event2 := <-progressStream
	assert.NoError(t, event2.Err)
}

func TestExecuteCdkCommandAndCleanupDirectory_Success(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	tempDirectory := "/my/directory"

	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osRemoveAll = mockOs.RemoveAll
	mockOs.EXPECT().RemoveAll(tempDirectory).Return(nil).Times(1)

	progressStream, _ := executeCdkCommandAndCleanupDirectory(t.TempDir(), []string{testExecuteCommandSuccessArg}, tempDirectory)
	event1 := <-progressStream
	assert.Equal(t, 3, event1.CurrentStep)
	assert.Equal(t, 10, event1.TotalSteps)
	assert.Equal(t, testExecuteCommandProgressLine, event1.Outputs[0])
	event2 := <-progressStream
	assert.NoError(t, event2.Err)
}

func TestExecuteCdkCommand_Failure(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	ctrl := gomock.NewController(t)
	mockOs := iomocks.NewMockOS(ctrl)
	osRemoveAll = mockOs.RemoveAll
	mockOs.EXPECT().RemoveAll(gomock.Any()).Return(nil).Times(0)

	progressStream, _ := executeCdkCommand(t.TempDir(), []string{testExecuteCommandFailureArg})
	event1 := <-progressStream
	assert.Equal(t, testExecuteCommandFailureArg, event1.Outputs[0])
	event2 := <-progressStream
	assert.Error(t, event2.Err)
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	require.GreaterOrEqual(t, len(args), 5)
	assert.Equal(t, "npm", args[0])
	assert.Equal(t, "run", args[1])
	assert.Equal(t, "cdk", args[2])
	assert.Equal(t, "--", args[3])

	testArg := args[4]
	switch testArg {
	case testExecuteCommandSuccessArg:
		fmt.Fprint(os.Stdout, "some line")
		fmt.Fprint(os.Stderr, testExecuteCommandProgressLine)
		os.Exit(0)
	case testExecuteCommandFailureArg:
		fmt.Fprint(os.Stdout, "some line")
		fmt.Fprint(os.Stderr, testExecuteCommandFailureArg)
		os.Exit(1)
	default:
		fmt.Fprint(os.Stderr, "Unknown failure")
		os.Exit(1)
	}
}
