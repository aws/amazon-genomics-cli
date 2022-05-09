package environment

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	"github.com/aws/amazon-genomics-cli/internal/pkg/constants"
)

var AllEngines = []string{
	constants.CROMWELL,
	constants.NEXTFLOW,
	constants.MINIWDL,
	constants.SNAKEMAKE,
	constants.TOIL,
}

var AllComponents []string = append(AllEngines, constants.WES)

var UsesWesAdapter = map[string]bool{
	constants.CROMWELL:  true,
	constants.NEXTFLOW:  true,
	constants.MINIWDL:   true,
	constants.SNAKEMAKE: true,
	constants.TOIL:      false,
}

const DefaultEcrRegistry = "680431765560"
const DefaultEcrRegion = "us-east-1"

var DefaultRepositories = map[string]string{
	constants.CROMWELL:  "aws/cromwell-mirror",
	constants.MINIWDL:   "aws/miniwdl-mirror",
	constants.NEXTFLOW:  "aws/nextflow-mirror",
	constants.SNAKEMAKE: "aws/snakemake-mirror",
	constants.TOIL:      "aws/toil-mirror",
	constants.WES:       "aws/wes-release",
}

// TODO: Implement better tag versioning system
var DefaultTags = map[string]string{
	constants.CROMWELL:  "64",
	constants.MINIWDL:   "v0.1.11",
	constants.NEXTFLOW:  "21.04.3",
	constants.SNAKEMAKE: "internal-fork",
	constants.TOIL:      "5.7.0a1-f77554b28bbda7b112c5b058b31d5041299c38c9",
	constants.WES:       "0.1.0",
}

func constructCommonImages() map[string]ecr.ImageReference {
	result := make(map[string]ecr.ImageReference)

	for _, component := range AllComponents {
		capName := strings.ToUpper(component)
		result[component] = ecr.ImageReference{
			RegistryId:     LookUpEnvOrDefault(fmt.Sprintf("ECR_%s_ACCOUNT_ID", capName), DefaultEcrRegistry),
			Region:         LookUpEnvOrDefault(fmt.Sprintf("ECR_%s_REGION", capName), DefaultEcrRegion),
			RepositoryName: LookUpEnvOrDefault(fmt.Sprintf("ECR_%s_REPOSITORY", capName), DefaultRepositories[component]),
			ImageTag:       LookUpEnvOrDefault(fmt.Sprintf("ECR_%s_TAG", capName), DefaultTags[component]),
		}
	}

	return result
}

var CommonImages = constructCommonImages()

func LookUpEnvOrDefault(envVariableName string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVariableName); ok {
		return value
	}
	return defaultValue
}
