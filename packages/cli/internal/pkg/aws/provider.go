package aws

import (
	"context"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/batch"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cfn"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cwl"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ddb"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ssm"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/sts"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/rs/zerolog/log"
)

type client string

const (
	clientCdk   client = "CDK"
	clientCfn   client = "CFN"
	clientCwl   client = "CWL"
	clientS3    client = "S3"
	clientSsm   client = "SSM"
	clientSts   client = "STS"
	clientDdb   client = "DDB"
	clientBatch client = "BATCH"
	clientEcr   client = "ECR"
)

var (
	profileConfigs = make(map[string]aws.Config)
	profileClients = make(map[string]map[client]interface{})
	loadConfig     = config.LoadDefaultConfig
)

func CdkClient(profile string) *cdk.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientCdk]; !ok {
		profileClients[profile][clientCdk] = cdk.NewClient(profile)
	}

	client := profileClients[profile][clientCdk].(cdk.Client)
	return &client
}

func CfnClient(profile string) *cfn.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientCfn]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientCfn] = cfn.New(cfg)
	}

	return profileClients[profile][clientCfn].(*cfn.Client)
}

func CwlClient(profile string) *cwl.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientCwl]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientCwl] = cwl.New(cfg)
	}

	return profileClients[profile][clientCwl].(*cwl.Client)
}

func S3Client(profile string) *s3.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientS3]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientS3] = s3.New(cfg)
	}

	return profileClients[profile][clientS3].(*s3.Client)
}

func SsmClient(profile string) *ssm.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientSsm]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientSsm] = ssm.New(cfg)
	}

	return profileClients[profile][clientSsm].(*ssm.Client)
}

func StsClient(profile string) *sts.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientSts]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientSts] = sts.NewClient(cfg)
	}

	client := profileClients[profile][clientSts].(sts.Client)
	return &client
}

func DdbClient(profile string) *ddb.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientDdb]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientDdb] = ddb.New(cfg)
	}

	client := profileClients[profile][clientDdb].(*ddb.Client)
	return client
}

func BatchClient(profile string) *batch.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientBatch]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientBatch] = batch.New(cfg)
	}

	client := profileClients[profile][clientBatch].(*batch.Client)
	return client
}

func EcrClient(profile string) *ecr.Client {
	initClientMap(profile)
	if _, ok := profileClients[profile][clientEcr]; !ok {
		cfg := GetProfileConfig(profile)
		profileClients[profile][clientEcr] = ecr.New(cfg)
	}

	client := profileClients[profile][clientEcr].(*ecr.Client)
	return client
}

func Region(profile string) string {
	initClientMap(profile)
	cfg := GetProfileConfig(profile)
	return cfg.Region
}

func initClientMap(profile string) {
	if _, ok := profileClients[profile]; !ok {
		profileClients[profile] = make(map[client]interface{})
	}
}

func GetProfileConfig(profile string) aws.Config {
	if _, ok := profileConfigs[profile]; !ok {
		cfg, err := loadConfig(context.Background(),
			config.WithSharedConfigProfile(profile),
			config.WithAssumeRoleCredentialOptions(func(o *stscreds.AssumeRoleOptions) {
				o.TokenProvider = stscreds.StdinTokenProvider
			}),
		)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		profileConfigs[profile] = cfg
	}

	return profileConfigs[profile]
}
