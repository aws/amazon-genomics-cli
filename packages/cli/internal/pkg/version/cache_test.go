package version

import (
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite

	cacheTimestamp    time.Time
	homeDir           string
	dataJson          string
	dataStruct        []Info
	expectedCachePath string
	version           string

	origUserHomeDir func() (string, error)
	origReadFile    func(filename string) ([]byte, error)
	origWriteFile   func(filename string, data []byte, perm fs.FileMode) error
}

func (s *CacheTestSuite) SetupTest() {
	s.cacheTimestamp, _ = time.Parse(time.RFC3339, "2021-10-11T23:50:50.104404138-07:00")
	s.homeDir = "/my/home/dir"
	s.expectedCachePath = "/my/home/dir/.agc/.version_cache"
	s.dataJson = `{
  "Versions": [
    {
      "Name": "1.0.1",
      "Deprecated": true,
      "DeprecationMessage": "Deprecation message test",
      "Highlight": "Highlight test"
    }
  ],
  "CurrentVersion": "1.0.1",
  "Timestamp": "2021-10-11T23:50:50.104404138-07:00"
}`
	s.dataStruct = []Info{
		{
			Name:               "1.0.1",
			Deprecated:         true,
			DeprecationMessage: "Deprecation message test",
			Highlight:          "Highlight test",
		},
	}
	s.version = "1.0.1"

	s.origUserHomeDir = userHomeDir
	s.origReadFile = readFile
	s.origWriteFile = writeFile

	userHomeDir = func() (string, error) {
		return s.homeDir, nil
	}

	readFile = func(filename string) ([]byte, error) {
		s.Assert().Equal(s.expectedCachePath, filename)
		return []byte(s.dataJson), nil
	}
}

func (s *CacheTestSuite) TearDownTest() {
	userHomeDir = s.origUserHomeDir
	readFile = s.origReadFile
	writeFile = s.origWriteFile
}

func (s *CacheTestSuite) TestReadFromCacheNominal() {
	expected := s.dataStruct
	now := s.cacheTimestamp.Add(1 * time.Hour)

	actual, err := readFromCache(s.version, now)
	if s.Assert().NoError(err) {
		s.Assert().Equal(expected, actual)
	}
}

func (s *CacheTestSuite) TestReadFromCacheExpired() {
	now := s.cacheTimestamp.Add(25 * time.Hour)
	_, err := readFromCache(s.version, now)
	if s.Assert().Error(err) {
		s.Assert().Equal(CacheExpiredError, err)
	}
}

func (s *CacheTestSuite) TestReadFromCacheNoUserHome() {
	errorMessage := "no user home available"
	userHomeDir = func() (string, error) {
		return "", errors.New(errorMessage)
	}
	readFile = func(filename string) ([]byte, error) {
		s.Fail("should not call 'readFile'")
		return nil, nil
	}
	_, err := readFromCache(s.version, s.cacheTimestamp)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
	}
}

func (s *CacheTestSuite) TestReadFromCacheNoFile() {
	errorMessage := "file not found"
	readFile = func(filename string) ([]byte, error) {
		return nil, errors.New(errorMessage)
	}
	_, err := readFromCache(s.version, s.cacheTimestamp)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
	}
}

func (s *CacheTestSuite) TestReadFromCacheInvalidFormat() {
	readFile = func(filename string) ([]byte, error) {
		s.Assert().Equal(s.expectedCachePath, filename)
		return []byte("<XML/>"), nil
	}
	_, err := readFromCache(s.version, s.cacheTimestamp)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, "invalid character '<' looking for beginning of value")
	}
}

func (s *CacheTestSuite) TestWriteToCacheNominal() {
	var actuallyWritten string
	writeFile = func(filename string, data []byte, perm fs.FileMode) error {
		s.Assert().Equal(s.expectedCachePath, filename)
		s.Assert().Equal(fs.FileMode(0644), perm)
		actuallyWritten = string(data)
		return nil
	}
	expected := `{"CurrentVersion":"1.0.1","Versions":[{"Name":"1.0.1","Deprecated":true,"DeprecationMessage":"Deprecation message test","Highlight":"Highlight test"}],"Timestamp":"2021-10-11T23:50:50.104404138-07:00"}`
	err := writeToCache(s.version, s.dataStruct, s.cacheTimestamp)
	if s.NoError(err) {
		s.Assert().Equal(expected, actuallyWritten)
	}
}

func (s *CacheTestSuite) TestWriteToCacheNoUserHome() {
	errorMessage := "no user home available"
	userHomeDir = func() (string, error) {
		return "", errors.New(errorMessage)
	}
	writeFile = func(filename string, data []byte, perm fs.FileMode) error {
		s.Fail("should not call 'readFile'")
		return nil
	}
	err := writeToCache(s.version, s.dataStruct, s.cacheTimestamp)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
	}
}

func (s *CacheTestSuite) TestWriteToCacheWriteFailure() {
	errorMessage := "cannot write file"
	writeFile = func(filename string, data []byte, perm fs.FileMode) error {
		return errors.New(errorMessage)
	}
	err := writeToCache(s.version, s.dataStruct, s.cacheTimestamp)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
	}
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
