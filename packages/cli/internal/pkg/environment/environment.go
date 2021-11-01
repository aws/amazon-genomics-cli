package environment

import (
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
)

const DefaultEcrRegistry = "680431765560"
const DefaultEcrRegion = "us-east-1"

// TODO: Implement better tag versioning system
const DefaultCromwellTag = "64"
const DefaultNextflowTag = "21.04.3"
const DefaultWesTag = "0.1.0"

const WesImageKey = "WES"
const CromwellImageKey = "CROMWELL"
const NextflowImageKey = "NEXTFLOW"

var CommonImages = map[string]ecr.ImageReference{
	WesImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_WES_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_WES_REGION", DefaultEcrRegion),
		RepositoryName: "aws/wes-release",
		ImageTag:       LookUpEnvOrDefault("ECR_WES_TAG", DefaultWesTag),
	},
	CromwellImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_CROMWELL_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_CROMWELL_REGION", DefaultEcrRegion),
		RepositoryName: "aws/cromwell-mirror",
		ImageTag:       LookUpEnvOrDefault("ECR_CROMWELL_TAG", DefaultCromwellTag),
	},
	NextflowImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_NEXTFLOW_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_NEXTFLOW_REGION", DefaultEcrRegion),
		RepositoryName: "aws/nextflow-mirror",
		ImageTag:       LookUpEnvOrDefault("ECR_NEXTFLOW_TAG", DefaultNextflowTag),
	},
}

func LookUpEnvOrDefault(envVariableName string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVariableName); ok {
		return value
	}
	return defaultValue
}
