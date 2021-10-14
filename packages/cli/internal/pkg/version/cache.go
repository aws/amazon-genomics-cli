package version

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheFileNameRel = ".agc/.version_cache"
)

var (
	CacheExpiredError = errors.New("cache expired")
	expirationTime    = 24 * time.Hour

	userHomeDir = os.UserHomeDir
	readFile    = ioutil.ReadFile
	writeFile   = ioutil.WriteFile
)

type CacheRecord struct {
	Versions  []Info
	Timestamp time.Time
}

var readFromCache = func(_ string, currentTime time.Time) ([]Info, error) {
	path, err := userHomeDir()
	if err != nil {
		return nil, err
	}
	cacheFileNameAbs := filepath.Join(path, cacheFileNameRel)
	cacheBites, err := readFile(cacheFileNameAbs)
	if err != nil {
		return nil, err
	}
	var record CacheRecord
	err = json.Unmarshal(cacheBites, &record)
	if err != nil {
		return nil, err
	}
	cacheAge := currentTime.Sub(record.Timestamp)
	if cacheAge > expirationTime {
		return nil, CacheExpiredError
	}
	return record.Versions, nil
}

var writeToCache = func(infos []Info, currentTime time.Time) error {
	path, err := userHomeDir()
	if err != nil {
		return err
	}
	cacheFileNameAbs := filepath.Join(path, cacheFileNameRel)

	cacheBites, err := json.Marshal(CacheRecord{infos, currentTime})
	if err != nil {
		return err
	}
	return writeFile(cacheFileNameAbs, cacheBites, 0644)
}
