package awsmocks

import (
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/batch"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ssm"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/sts"
)

type CdkClient interface {
	cdk.Interface
	DisplayProgressBar(description string, progressEvents []cdk.ProgressStream) []cdk.Result
	ShowExecution(progressEvents []cdk.ProgressStream) []cdk.Result
	SilentExecution(progressStreams []cdk.ProgressStream) []cdk.Result
}

type S3Client interface {
	s3.Interface
}

type StsClient interface {
	sts.Interface
}

type SsmClient interface {
	ssm.Interface
}

type CfnClient interface {
	cfn.Interface
}

type BatchClient interface {
	batch.Interface
}

type CwlClient interface {
	cwl.Interface
}

type CwlLogPaginator interface {
	cwl.LogPaginator
}

type DdbClient interface {
	ddb.Interface
}

type EcrClient interface {
	ecr.Interface
}
