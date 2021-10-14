package version

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite

	version   string
	timestamp time.Time
	infos     []Info

	origReadFromCache func(version string, currentTime time.Time) ([]Info, error)
	origWriteToCache  func(version string, infos []Info, currentTime time.Time) error
}

func (s *StoreTestSuite) SetupTest() {
	s.version = "1.0.1"
	s.timestamp, _ = time.Parse(time.RFC3339, "2021-10-11T15:04:05Z07:00")
	s.infos = []Info{{Name: s.version}}

	s.origReadFromCache = readFromCache
	s.origWriteToCache = writeToCache
}

func (s *StoreTestSuite) TearDownTest() {
	readFromCache = s.origReadFromCache
	writeToCache = s.origWriteToCache
}

func (s *StoreTestSuite) TestStoreReadFromCache() {
	readFromCache = func(_ string, _ time.Time) ([]Info, error) {
		return s.infos, nil
	}

	store := &cachedStore{
		ReplenishFunc: func(version string) ([]Info, error) {
			s.Assert().Equal(s.version, version)
			s.Fail("should not replenish the cache")
			return nil, nil
		},
	}

	writeToCache = func(_ string, _ []Info, _ time.Time) error {
		s.Fail("should overwrite the cache")
		return nil
	}

	actual, err := store.ReadVersions(s.version, s.timestamp)
	if s.Assert().NoError(err) {
		s.Assert().Equal(s.infos, actual)
	}
}

func (s *StoreTestSuite) TestStoreExpiredCache() {
	readFromCache = func(_ string, _ time.Time) ([]Info, error) {
		return nil, CacheExpiredError
	}

	store := &cachedStore{
		ReplenishFunc: func(version string) ([]Info, error) {
			s.Assert().Equal(s.version, version)
			return s.infos, nil
		},
	}

	writeToCache = func(version string, infos []Info, timestamp time.Time) error {
		s.Assert().Equal(s.version, version)
		s.Assert().Equal(s.infos, infos)
		s.Assert().Equal(s.timestamp, timestamp)
		return nil
	}

	actual, err := store.ReadVersions(s.version, s.timestamp)
	if s.Assert().NoError(err) {
		s.Assert().Equal(s.infos, actual)
	}
}

func (s *StoreTestSuite) TestStoreReplenishFails() {
	errorMessage := "failed to replenish"
	readFromCache = func(_ string, _ time.Time) ([]Info, error) {
		return nil, CacheExpiredError
	}

	store := &cachedStore{
		ReplenishFunc: func(_ string) ([]Info, error) {
			return s.infos, errors.New(errorMessage)
		},
	}

	writeToCache = func(_ string, _ []Info, _ time.Time) error {
		s.Fail("should not write to cache")
		return nil
	}

	_, err := store.ReadVersions(s.version, s.timestamp)
	if s.Assert().Error(err) {
		s.Assert().EqualError(err, errorMessage)
	}
}

func (s *StoreTestSuite) TestStoreCacheUpdateFails() {
	errorMessage := "failed to update the cache"
	readFromCache = func(_ string, _ time.Time) ([]Info, error) {
		return nil, CacheExpiredError
	}

	store := &cachedStore{
		ReplenishFunc: func(version string) ([]Info, error) {
			s.Assert().Equal(s.version, version)
			return s.infos, nil
		},
	}

	writeToCache = func(_ string, _ []Info, _ time.Time) error {
		return errors.New(errorMessage)
	}

	actual, err := store.ReadVersions(s.version, s.timestamp)
	if s.Assert().NoError(err) {
		s.Assert().Equal(s.infos, actual)
	}
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
