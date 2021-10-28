package version

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/suite"
)

type CheckTestSuite struct {
	suite.Suite

	now time.Time

	origNewS3ClientFromConfig func(cfg aws.Config) S3Api
	origCheckVersion          func(s3Reader S3Api, channel string, currentTime time.Time) (Result, error)
	origGetCurrentTime        func() time.Time
}

func (suite *CheckTestSuite) SetupTest() {
	suite.now, _ = time.Parse(time.RFC3339, "2021-10-10T15:04:05Z07:00")
	Version = "1.0.0-41-gc5ac696"

	suite.origNewS3ClientFromConfig = newS3ClientFromConfig
	suite.origCheckVersion = checkVersion
	suite.origGetCurrentTime = getCurrentTime
}

func (suite *CheckTestSuite) TearDownTest() {
	newS3ClientFromConfig = suite.origNewS3ClientFromConfig
	checkVersion = suite.origCheckVersion
	getCurrentTime = suite.origGetCurrentTime
}

func (suite *CheckTestSuite) TestCheckDefaultChannel() {
	getCurrentTime = func() time.Time {
		return suite.now
	}

	testResult := Result{CurrentVersion: Version}

	checkVersion = func(s3Reader S3Api, channel string, currentTime time.Time) (Result, error) {
		suite.Assert().Equal(channel, DefaultChannel)
		suite.Assert().Equal(currentTime, suite.now)
		return testResult, nil
	}

	newS3ClientFromConfig = func(cfg aws.Config) S3Api {
		suite.Assert().Equal(cfg.Region, "us-east-1")
		return nil
	}

	result, err := Check()
	if suite.Assert().NoError(err) {
		suite.Assert().Equal(Version, result.CurrentVersion)
	}
}

func (suite *CheckTestSuite) TestCheckDefaultCustomChannel() {
	getCurrentTime = func() time.Time {
		return suite.now
	}

	testResult := Result{CurrentVersion: Version}

	const customChannel = "s3://custom"
	suite.T().Setenv(ChannelVarName, customChannel)

	checkVersion = func(s3Reader S3Api, channel string, currentTime time.Time) (Result, error) {
		suite.Assert().Equal(channel, customChannel)
		suite.Assert().Equal(currentTime, suite.now)
		return testResult, nil
	}

	newS3ClientFromConfig = func(cfg aws.Config) S3Api {
		suite.Assert().Equal(cfg.Region, "us-east-1")
		return nil
	}

	result, err := Check()
	if suite.Assert().NoError(err) {
		suite.Assert().Equal(Version, result.CurrentVersion)
	}
}

func (suite *CheckTestSuite) TestCheckSkip() {
	getCurrentTime = func() time.Time {
		return suite.now
	}

	suite.T().Setenv(UpdateCheckCtrlVarName, "false")

	checkVersion = func(s3Reader S3Api, channel string, currentTime time.Time) (Result, error) {
		suite.Fail("Should not call 'checkVersion'")
		return Result{}, nil
	}

	newS3ClientFromConfig = func(cfg aws.Config) S3Api {
		suite.Fail("Should not call 'newS3ClientFromConfig'")
		return nil
	}

	result, err := Check()
	if suite.Assert().NoError(err) {
		suite.Assert().Equal(Version, result.CurrentVersion)
	}
}

func TestCheckTestSuite(t *testing.T) {
	suite.Run(t, new(CheckTestSuite))
}
