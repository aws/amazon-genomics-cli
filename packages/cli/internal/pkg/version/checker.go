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
	result := Result{CurrentVersion: currentVersion}
	infos, err := ch.store.ReadVersions(currentVersion, ch.currentTime)
	if err != nil {
		return Result{}, err
	}
	if len(infos) == 0 {
		return Result{}, nil
	}
	if infos[0].Name == currentVersion {
		currentInfo := infos[0]
		infos = infos[1:]
		result.CurrentVersionDeprecated = currentInfo.Deprecated
		result.CurrentVersionDeprecationMessage = currentInfo.DeprecationMessage
	}
	if len(infos) == 0 {
		result.LatestVersion = result.CurrentVersion
		return result, nil
	}

	result.LatestVersion = infos[len(infos)-1].Name
	for _, info := range infos {
		if len(info.Highlight) > 0 {
			result.NewerVersionHighlights = append(result.NewerVersionHighlights, info.Highlight)
		}
	}
	return result, nil
}

func finalizeVersion(v *semver.Version) {
	v.Pre = nil
	v.Build = nil
}
