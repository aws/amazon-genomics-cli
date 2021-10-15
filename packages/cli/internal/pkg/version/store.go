package version

import (
	"time"

	"github.com/rs/zerolog/log"
)

type cachedStore struct {
	ReplenishFunc func(versionString string) ([]Info, error)
}

func (s *cachedStore) ReadVersions(version string, currentTime time.Time) ([]Info, error) {
	infos, err := readFromCache(version, currentTime)
	if err != nil {
		log.Debug().Msgf("failed to read local cache: %v", err)
		infos, err = s.ReplenishFunc(version)
		if err != nil {
			return nil, err
		}
		err = writeToCache(version, infos, currentTime)
		if err != nil {
			log.Debug().Msgf("failed to write local cache %v", err)
		}
	}
	return infos, nil
}
