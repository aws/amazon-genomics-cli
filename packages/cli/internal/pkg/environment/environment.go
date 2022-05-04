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
const DefaultMiniwdlTag = "v0.1.11"
const DefaultSnakemakeTag = "internal-fork"
const DefaultToilTag = "5.7.0a1-f77554b28bbda7b112c5b058b31d5041299c38c9"

const WesImageKey = "WES"
const CromwellImageKey = "CROMWELL"
const NextflowImageKey = "NEXTFLOW"
const MiniwdlImageKey = "MINIWDL"
const SnakemakeImageKey = "SNAKEMAKE"
const ToilImageKey = "TOIL"

var CommonImages = map[string]ecr.ImageReference{
	WesImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_WES_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_WES_REGION", DefaultEcrRegion),
		RepositoryName: LookUpEnvOrDefault("ECR_WES_REPOSITORY", "aws/wes-release"),
		ImageTag:       LookUpEnvOrDefault("ECR_WES_TAG", DefaultWesTag),
	},
	CromwellImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_CROMWELL_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_CROMWELL_REGION", DefaultEcrRegion),
		RepositoryName: LookUpEnvOrDefault("ECR_CROMWELL_REPOSITORY", "aws/cromwell-mirror"),
		ImageTag:       LookUpEnvOrDefault("ECR_CROMWELL_TAG", DefaultCromwellTag),
	},
	NextflowImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_NEXTFLOW_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_NEXTFLOW_REGION", DefaultEcrRegion),
		RepositoryName: LookUpEnvOrDefault("ECR_NEXTFLOW_REPOSITORY", "aws/nextflow-mirror"),
		ImageTag:       LookUpEnvOrDefault("ECR_NEXTFLOW_TAG", DefaultNextflowTag),
	},
	MiniwdlImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_MINIWDL_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_MINIWDL_REGION", DefaultEcrRegion),
		RepositoryName: LookUpEnvOrDefault("ECR_MINIWDL_REPOSITORY", "aws/miniwdl-mirror"),
		ImageTag:       LookUpEnvOrDefault("ECR_MINIWDL_TAG", DefaultMiniwdlTag),
	},
	SnakemakeImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_SNAKEMAKE_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_SNAKEMAKE_REGION", DefaultEcrRegion),
		RepositoryName: LookUpEnvOrDefault("ECR_SNAKEMAKE_REPOSITORY", "aws/snakemake-mirror"),
		ImageTag:       LookUpEnvOrDefault("ECR_SNAKEMAKE_TAG", DefaultSnakemakeTag),
	},
	ToilImageKey: {
		RegistryId:     LookUpEnvOrDefault("ECR_TOIL_ACCOUNT_ID", DefaultEcrRegistry),
		Region:         LookUpEnvOrDefault("ECR_TOIL_REGION", DefaultEcrRegion),
		RepositoryName: LookUpEnvOrDefault("ECR_TOIL_REPOSITORY", "aws/toil-mirror"),
		ImageTag:       LookUpEnvOrDefault("ECR_TOIL_TAG", DefaultToilTag),
	},
}

// Some workflow engines require other images
var ImageDependencies = map[string]([]string){
	WesImageKey:       {},
	CromwellImageKey:  {WesImageKey},
	NextflowImageKey:  {WesImageKey},
	MiniwdlImageKey:   {WesImageKey},
	SnakemakeImageKey: {WesImageKey},
	ToilImageKey:      {},
}

func LookUpEnvOrDefault(envVariableName string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVariableName); ok {
		return value
	}
	return defaultValue
}
