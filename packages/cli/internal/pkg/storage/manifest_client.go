package storage

import (
	"fmt"
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/spec"
)

const ManifestFileName = "MANIFEST.json"

var specFromJson = spec.FromJson
var osStat = os.Stat

func ReadManifestInDirectory(directory string) (spec.Manifest, error) {
	return specFromJson(getManifestFilepath(directory))
}

func DoesManifestExistInDirectory(directory string) bool {
	if _, err := osStat(getManifestFilepath(directory)); err != nil {
		return false
	}
	return true
}

func getManifestFilepath(directory string) string {
	return fmt.Sprintf("%s/%s", directory, ManifestFileName)
}
