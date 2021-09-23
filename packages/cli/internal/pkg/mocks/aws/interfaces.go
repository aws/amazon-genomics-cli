package awsmocks

import (
	"github.com/aws/amazon-genomics-cli/common/aws/batch"
	"github.com/aws/amazon-genomics-cli/common/aws/cdk"
	"github.com/aws/amazon-genomics-cli/common/aws/cfn"
	"github.com/aws/amazon-genomics-cli/common/aws/cwl"
	"github.com/aws/amazon-genomics-cli/common/aws/ddb"
	"github.com/aws/amazon-genomics-cli/common/aws/ecr"
	"github.com/aws/amazon-genomics-cli/common/aws/s3"
	"github.com/aws/amazon-genomics-cli/common/aws/ssm"
	"github.com/aws/amazon-genomics-cli/common/aws/sts"
)

type CdkClient interface {
	cdk.Interface
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
