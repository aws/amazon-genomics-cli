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

	origUserHomeDir func() (string, error)
	origReadFile    func(filename string) ([]byte, error)
	origWriteFile   func(filename string, data []byte, perm fs.FileMode) error
}

func (suite *CacheTestSuite) SetupTest() {
	suite.cacheTimestamp, _ = time.Parse(time.RFC3339, "2021-10-11T23:50:50.104404138-07:00")
	suite.homeDir = "/my/home/dir"
	suite.expectedCachePath = "/my/home/dir/.agc/.version_cache"
	suite.dataJson = `{
  "Versions": [
    {
      "Name": "1.0.1",
      "Deprecated": true,
      "DeprecationMessage": "Deprecation message test",
      "Highlight": "Highlight test"
    }
  ],
  "Timestamp": "2021-10-11T23:50:50.104404138-07:00"
}`
	suite.dataStruct = []Info{
		{
			Name:               "1.0.1",
			Deprecated:         true,
			DeprecationMessage: "Deprecation message test",
			Highlight:          "Highlight test",
		},
	}

	suite.origUserHomeDir = userHomeDir
	suite.origReadFile = readFile
	suite.origWriteFile = writeFile

	userHomeDir = func() (string, error) {
		return suite.homeDir, nil
	}

	readFile = func(filename string) ([]byte, error) {
		suite.Assert().Equal(suite.expectedCachePath, filename)
		return []byte(suite.dataJson), nil
	}
}

func (suite *CacheTestSuite) TearDownTest() {
	userHomeDir = suite.origUserHomeDir
	readFile = suite.origReadFile
	writeFile = suite.origWriteFile
}

func (suite *CacheTestSuite) TestReadFromCacheNominal() {
	expected := suite.dataStruct
	now := suite.cacheTimestamp.Add(1 * time.Hour)

	actual, err := readFromCache("", now)
	if suite.Assert().NoError(err) {
		suite.Assert().Equal(expected, actual)
	}
}

func (suite *CacheTestSuite) TestReadFromCacheExpired() {
	now := suite.cacheTimestamp.Add(25 * time.Hour)
	_, err := readFromCache("", now)
	if suite.Assert().Error(err) {
		suite.Assert().Equal(CacheExpiredError, err)
	}
}

func (suite *CacheTestSuite) TestReadFromCacheNoUserHome() {
	errorMessage := "no user home available"
	userHomeDir = func() (string, error) {
		return "", errors.New(errorMessage)
	}
	readFile = func(filename string) ([]byte, error) {
		suite.Fail("should not call 'readFile'")
		return nil, nil
	}
	_, err := readFromCache("", suite.cacheTimestamp)
	if suite.Assert().Error(err) {
		suite.Assert().EqualError(err, errorMessage)
	}
}

func (suite *CacheTestSuite) TestReadFromCacheNoFile() {
	errorMessage := "file not found"
	readFile = func(filename string) ([]byte, error) {
		return nil, errors.New(errorMessage)
	}
	_, err := readFromCache("", suite.cacheTimestamp)
	if suite.Assert().Error(err) {
		suite.Assert().EqualError(err, errorMessage)
	}
}

func (suite *CacheTestSuite) TestReadFromCacheInvalidFormat() {
	readFile = func(filename string) ([]byte, error) {
		suite.Assert().Equal(suite.expectedCachePath, filename)
		return []byte("<XML/>"), nil
	}
	_, err := readFromCache("", suite.cacheTimestamp)
	if suite.Assert().Error(err) {
		suite.Assert().EqualError(err, "invalid character '<' looking for beginning of value")
	}
}

func (suite *CacheTestSuite) TestWriteToCacheNominal() {
	var actuallyWritten string
	writeFile = func(filename string, data []byte, perm fs.FileMode) error {
		suite.Assert().Equal(suite.expectedCachePath, filename)
		suite.Assert().Equal(fs.FileMode(0644), perm)
		actuallyWritten = string(data)
		return nil
	}
	expected := `{"Versions":[{"Name":"1.0.1","Deprecated":true,"DeprecationMessage":"Deprecation message test","Highlight":"Highlight test"}],"Timestamp":"2021-10-11T23:50:50.104404138-07:00"}`
	err := writeToCache(suite.dataStruct, suite.cacheTimestamp)
	if suite.NoError(err) {
		suite.Assert().Equal(expected, actuallyWritten)
	}
}

func (suite *CacheTestSuite) TestWriteToCacheNoUserHome() {
	errorMessage := "no user home available"
	userHomeDir = func() (string, error) {
		return "", errors.New(errorMessage)
	}
	writeFile = func(filename string, data []byte, perm fs.FileMode) error {
		suite.Fail("should not call 'readFile'")
		return nil
	}
	err := writeToCache(suite.dataStruct, suite.cacheTimestamp)
	if suite.Assert().Error(err) {
		suite.Assert().EqualError(err, errorMessage)
	}
}

func (suite *CacheTestSuite) TestWriteToCacheWriteFailure() {
	errorMessage := "cannot write file"
	writeFile = func(filename string, data []byte, perm fs.FileMode) error {
		return errors.New(errorMessage)
	}
	err := writeToCache(suite.dataStruct, suite.cacheTimestamp)
	if suite.Assert().Error(err) {
		suite.Assert().EqualError(err, errorMessage)
	}
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
