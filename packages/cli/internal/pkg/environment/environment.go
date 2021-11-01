package environment

import (
	"os"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
)

const DefaultEcrRegistry = "555741984805"
const DefaultEcrRegion = "us-east-1"

const WesImageKey = "WES"
const CromwellImageKey = "CROMWELL"
const NextflowImageKey = "NEXTFLOW"
const MiniwdlImageKey = "MINIWDL"

var CommonImages = map[string]ecr.ImageReference{
	WesImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_WES_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_WES_REGION", DefaultEcrRegion),
		RepositoryName: "agc-wes-adapter-cromwell",
		ImageTag:       LookUpEnvOrDefault("ECR_WES_TAG", "WES_ECR_TAG_PLACEHOLDER"),
	},
	CromwellImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_CROMWELL_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_CROMWELL_REGION", DefaultEcrRegion),
		RepositoryName: "cromwell",
		ImageTag:       LookUpEnvOrDefault("ECR_CROMWELL_TAG", "CROMWELL_ECR_TAG_PLACEHOLDER"),
	},
	NextflowImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_NEXTFLOW_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_NEXTFLOW_REGION", DefaultEcrRegion),
		RepositoryName: "nextflow",
		ImageTag:       LookUpEnvOrDefault("ECR_NEXTFLOW_TAG", "NEXTFLOW_ECR_TAG_PLACEHOLDER"),
	},
	MiniwdlImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_MINIWDL_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_MINIWDL_REGION", DefaultEcrRegion),
		RepositoryName: "miniwdl",
		ImageTag:       LookUpEnvOrDefault("ECR_MINIWDL_TAG", "MINIWDL_ECR_TAG_PLACEHOLDER"),
	},
}

func LookUpEnvOrDefault(envVariableName string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVariableName); ok {
		return value
	}
	return defaultValue
}
