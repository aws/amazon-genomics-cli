package version

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type CheckerCheckTestSuite struct {
	suite.Suite

	ctrl               *gomock.Controller
	mockStore          *MockStore
	testTime           time.Time
	currentVersion     string
	nextVersion        string
	latestVersion      string
	currentVersionInfo Info
	nextVersionInfo    Info
	latestVersionInfo  Info
}

func (s *CheckerCheckTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.mockStore = NewMockStore(s.ctrl)
	s.testTime, _ = time.Parse(time.RFC3339, "2021-10-11T15:04:05Z07:00")
	s.currentVersion = "1.0.1"
	s.nextVersion = "1.3.5"
	s.latestVersion = "2.0.0"
	s.currentVersionInfo = Info{Name: s.currentVersion}
	s.nextVersionInfo = Info{Name: s.nextVersion}
	s.latestVersionInfo = Info{Name: s.latestVersion}
}

func (s *CheckerCheckTestSuite) AfterTest(_, _ string) {
	s.ctrl.Finish()
}

func (s *CheckerCheckTestSuite) TestCheckNominal() {
	s.mockStore.EXPECT().ReadVersions(s.currentVersion, s.testTime).
		Return([]Info{s.currentVersionInfo, s.nextVersionInfo, s.latestVersionInfo}, nil)
	checker := &checker{s.mockStore, s.testTime}
	expected := Result{
		CurrentVersion: s.currentVersion,
		LatestVersion:  s.latestVersion,
	}
	actual, err := checker.Check(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *CheckerCheckTestSuite) TestCheckNoVersions() {
	s.mockStore.EXPECT().ReadVersions(s.currentVersion, s.testTime).
		Return([]Info{}, nil)
	checker := &checker{s.mockStore, s.testTime}
	expected := Result{CurrentVersion: s.currentVersion, LatestVersion: s.currentVersion}
	actual, err := checker.Check(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *CheckerCheckTestSuite) TestCheckOneVersion() {
	s.mockStore.EXPECT().ReadVersions(s.currentVersion, s.testTime).
		Return([]Info{s.currentVersionInfo}, nil)
	checker := &checker{s.mockStore, s.testTime}
	expected := Result{
		CurrentVersion: s.currentVersion,
		LatestVersion:  s.currentVersion,
	}
	actual, err := checker.Check(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *CheckerCheckTestSuite) TestCheckDeprecatedVersion() {
	s.currentVersionInfo.Deprecated = true
	s.currentVersionInfo.DeprecationMessage = "Deprecated!"
	s.mockStore.EXPECT().ReadVersions(s.currentVersion, s.testTime).
		Return([]Info{s.currentVersionInfo, s.nextVersionInfo, s.latestVersionInfo}, nil)
	checker := &checker{s.mockStore, s.testTime}
	expected := Result{
		CurrentVersion:                   s.currentVersion,
		LatestVersion:                    s.latestVersion,
		CurrentVersionDeprecated:         true,
		CurrentVersionDeprecationMessage: "Deprecated!",
	}
	actual, err := checker.Check(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *CheckerCheckTestSuite) TestCheckVersionsWithHighlights() {
	s.nextVersionInfo.Highlight = "Great version!"
	s.latestVersionInfo.Highlight = "Even better version!"
	s.mockStore.EXPECT().ReadVersions(s.currentVersion, s.testTime).
		Return([]Info{s.currentVersionInfo, s.nextVersionInfo, s.latestVersionInfo}, nil)
	checker := &checker{s.mockStore, s.testTime}
	expected := Result{
		CurrentVersion: s.currentVersion,
		LatestVersion:  s.latestVersion,
		NewerVersionHighlights: []string{
			"Great version!",
			"Even better version!",
		},
	}
	actual, err := checker.Check(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *CheckerCheckTestSuite) TestCheckMissingCurrentVersion() {
	s.mockStore.EXPECT().ReadVersions(s.currentVersion, s.testTime).
		Return([]Info{s.nextVersionInfo, s.latestVersionInfo}, nil)
	checker := &checker{s.mockStore, s.testTime}
	expected := Result{
		CurrentVersion: s.currentVersion,
		LatestVersion:  s.latestVersion,
	}
	actual, err := checker.Check(s.currentVersion)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *CheckerCheckTestSuite) TestCheckInvalidVersion() {
	checker := &checker{s.mockStore, s.testTime}
	_, err := checker.Check("foo")
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, "No Major.Minor.Patch elements found")
	}
}

func (s *CheckerCheckTestSuite) TestCheckFailedStore() {
	errorMessage := "cannot read versions"
	s.mockStore.EXPECT().ReadVersions(s.currentVersion, s.testTime).
		Return(nil, errors.New(errorMessage))
	checker := &checker{s.mockStore, s.testTime}
	_, err := checker.Check(s.currentVersion)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
	}
}

func TestCheckerCheckTestSuite(t *testing.T) {
	suite.Run(t, new(CheckerCheckTestSuite))
}
