package environment

import (
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
)

const DefaultEcrRegistry = "555741984805"
const DefaultEcrRegion = "us-east-1"

// TODO: Implement better tag versioning system
const DefaultCromwellTag = "2021-10-01T21-33-26Z"
const DefaultNextflowTag = "2021-10-01T21-33-26Z"
const DefaultWesTag = "2021-10-01T21-33-26Z"

const WesImageKey = "WES"
const CromwellImageKey = "CROMWELL"
const NextflowImageKey = "NEXTFLOW"

var CommonImages = map[string]ecr.ImageReference{
	WesImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_WES_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_WES_REGION", DefaultEcrRegion),
		RepositoryName: "agc-wes-adapter-cromwell",
		ImageTag:       LookUpEnvOrDefault("ECR_WES_TAG", DefaultWesTag),
	},
	CromwellImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_CROMWELL_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_CROMWELL_REGION", DefaultEcrRegion),
		RepositoryName: "cromwell",
		ImageTag:       LookUpEnvOrDefault("ECR_CROMWELL_TAG", DefaultCromwellTag),
	},
	NextflowImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_NEXTFLOW_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_NEXTFLOW_REGION", DefaultEcrRegion),
		RepositoryName: "nextflow",
		ImageTag:       LookUpEnvOrDefault("ECR_NEXTFLOW_TAG", DefaultNextflowTag),
	},
}

func LookUpEnvOrDefault(envVariableName string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVariableName); ok {
		return value
	}
	return defaultValue
}
