package version

import (
	"time"

	"github.com/blang/semver/v4"
)

type checker struct {
	store Store

	currentTime time.Time
}

func (ch *checker) Check(versionFullString string) (Result, error) {
	currentSemVersion, err := semver.Parse(versionFullString)
	if err != nil {
		return Result{}, err
	}
	finalizeVersion(&currentSemVersion)
	currentVersion := currentSemVersion.String()
	infos, err := ch.store.ReadVersions(currentVersion, ch.currentTime)
	if err != nil {
		return Result{}, err
	}
	result := Result{
		CurrentVersion: currentVersion,
		LatestVersion:  currentVersion,
	}
	for i, info := range infos {
		if i == 0 && currentVersion == info.Name {
			result.CurrentVersionDeprecated = info.Deprecated
			result.CurrentVersionDeprecationMessage = info.DeprecationMessage
			continue
		}
		if len(info.Highlight) > 0 {
			result.NewerVersionHighlights = append(result.NewerVersionHighlights, info.Highlight)
		}
		result.LatestVersion = info.Name
	}
	return result, nil
}

func finalizeVersion(v *semver.Version) {
	v.Pre = nil
	v.Build = nil
}
