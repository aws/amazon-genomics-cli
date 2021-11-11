package version

type Result struct {
	CurrentVersion                   string
	LatestVersion                    string
	CurrentVersionDeprecated         bool
	CurrentVersionDeprecationMessage string
	NewerVersionHighlights           []string
}
