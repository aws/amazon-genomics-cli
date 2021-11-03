package cdk

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	iomocks "github.com/aws/amazon-genomics-cli/internal/pkg/mocks/io"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

type ExecuteCdkCommandTestSuite struct {
	suite.Suite

	osRemoveAllOrig func(string) error
	execCommandOrig func(command string, args ...string) *exec.Cmd

	ctrl   *gomock.Controller
	mockOs *iomocks.MockOS
	appDir string
	tmpDir string
}

func TestExecuteCdkCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ExecuteCdkCommandTestSuite))
}

func (s *ExecuteCdkCommandTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockOs = iomocks.NewMockOS(s.ctrl)
	s.osRemoveAllOrig = osRemoveAll
	s.execCommandOrig = execCommand

	osRemoveAll = s.mockOs.RemoveAll
	execCommand = fakeExecCommand

	s.appDir = s.T().TempDir()
	s.tmpDir = "/test/tmp/dir"
}

func (s *ExecuteCdkCommandTestSuite) AfterTest(_, _ string) {
	s.ctrl.Finish()
}

func (s *ExecuteCdkCommandTestSuite) TearDownTest() {
	osRemoveAll = s.osRemoveAllOrig
	execCommand = s.execCommandOrig
}

func (s *ExecuteCdkCommandTestSuite) TestExecuteCdkCommand_Success() {
	s.mockOs.EXPECT().RemoveAll(gomock.Any()).Return(nil).Times(0)

	progressStream, err := executeCdkCommand(s.appDir, []string{testExecuteCommandSuccessArg})
	s.Require().NoError(err)
	event1 := <-progressStream
	s.Assert().Equal(3, event1.CurrentStep)
	s.Assert().Equal(10, event1.TotalSteps)
	s.Assert().Equal(testExecuteCommandProgressLine, event1.Outputs[0])
	event2 := <-progressStream
	s.Assert().NoError(event2.Err)
	waitForChanToClose(progressStream)
}

func (s *ExecuteCdkCommandTestSuite) TestExecuteCdkCommandAndCleanupDirectory_Success() {
	s.mockOs.EXPECT().RemoveAll(s.tmpDir).Return(nil).Times(1)

	progressStream, err := executeCdkCommandAndCleanupDirectory(s.appDir, []string{testExecuteCommandSuccessArg}, s.tmpDir)
	s.Require().NoError(err)
	event1 := <-progressStream
	s.Assert().Equal(3, event1.CurrentStep)
	s.Assert().Equal(10, event1.TotalSteps)
	s.Assert().Equal(testExecuteCommandProgressLine, event1.Outputs[0])
	event2 := <-progressStream
	s.Assert().NoError(event2.Err)
	waitForChanToClose(progressStream)
}

func (s *ExecuteCdkCommandTestSuite) TestExecuteCdkCommandAndCleanupDirectory_Failure() {
	s.mockOs.EXPECT().RemoveAll(s.tmpDir).Return(nil).Times(1)

	progressStream, err := executeCdkCommandAndCleanupDirectory(s.appDir, []string{testExecuteCommandFailureArg}, s.tmpDir)
	s.Require().NoError(err)
	event1 := <-progressStream
	s.Assert().Equal(testExecuteCommandFailureArg, event1.Outputs[0])
	event2 := <-progressStream
	s.Assert().Error(event2.Err)
	waitForChanToClose(progressStream)
}

func (s *ExecuteCdkCommandTestSuite) TestExecuteCdkCommandAndCleanupDirectory_FailToExecute() {
	s.mockOs.EXPECT().RemoveAll(s.tmpDir).Return(nil).Times(1)

	progressStream, err := executeCdkCommandAndCleanupDirectory("foo/bar", []string{testExecuteCommandFailureArg}, s.tmpDir)
	s.Assert().Error(err)
	s.Assert().Nil(progressStream)
}

func (s *ExecuteCdkCommandTestSuite) TestExecuteCdkCommand_Failure() {
	s.mockOs.EXPECT().RemoveAll(gomock.Any()).Return(nil).Times(0)

	progressStream, err := executeCdkCommand(s.appDir, []string{testExecuteCommandFailureArg})
	s.Require().NoError(err)
	event1 := <-progressStream
	s.Assert().Equal(testExecuteCommandFailureArg, event1.Outputs[0])
	event2 := <-progressStream
	s.Assert().Error(event2.Err)
	waitForChanToClose(progressStream)
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

func waitForChanToClose(channel ProgressStream) {
	for range channel {
	}
}
