// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror/actionableerror"
	"path/filepath"

	"github.com/aws/amazon-genomics-cli/internal/pkg/aws"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/cdk"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/ecr"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/s3"
	"github.com/aws/amazon-genomics-cli/internal/pkg/aws/sts"
	"github.com/aws/amazon-genomics-cli/internal/pkg/cli/clierror"
	"github.com/aws/amazon-genomics-cli/internal/pkg/logging"
	"github.com/aws/amazon-genomics-cli/internal/pkg/osutils"
	"github.com/aws/amazon-genomics-cli/internal/pkg/version"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	accountBucketFlag            = "bucket"
	accountVpcFlag               = "vpc"
	publicSubnetsFlag            = "usePublicSubnets"
	accountBucketFlagDescription = `The name of an S3 bucket that AGC will use to store its data.
An autogenerated name will be used if not specified. A new bucket will be created if the bucket does not exist.`
	accountVpcFlagDescription = `The ID of a VPC that AGC will run in. 
A new VPC will be created if not specified.`
	publicSubnetsFlagDescription = `Do not create a NAT gateway or VPC endpoints, to lower operating costs. 
You must enable the usePublicSubnets option in your project context if you use this flag.`
	cdkCoreDir   = ".agc/cdk/apps/core"
	bucketPrefix = "agc"
	activateKey  = "activate"
	bootstrapKey = "bootstrap"
)

type accountActivateVars struct {
	bucketName    string
	vpcId         string
	publicSubnets bool
}

type accountActivateOpts struct {
	accountActivateVars
	stsClient sts.Interface
	s3Client  s3.Interface
	cdkClient cdk.Interface
	ecrClient ecr.Interface
	imageRefs map[string]ecr.ImageReference
	region    string
}

func newAccountActivateOpts(vars accountActivateVars) (*accountActivateOpts, error) {
	return &accountActivateOpts{
		accountActivateVars: vars,
		stsClient:           aws.StsClient(profile),
		s3Client:            aws.S3Client(profile),
		cdkClient:           cdk.NewClient(profile),
		region:              aws.Region(profile),
	}, nil
}

// Execute activates AGC.
func (o *accountActivateOpts) Execute() error {
	if o.bucketName == "" {
		bucketName, err := o.generateDefaultBucket()
		if err != nil {
			return err
		}
		o.bucketName = bucketName
	}

	exists, err := o.s3Client.BucketExists(o.bucketName)
	if err != nil {
		return err
	}

	environmentVars := []string{
		fmt.Sprintf("AGC_BUCKET_NAME=%s", o.bucketName),
		fmt.Sprintf("CREATE_AGC_BUCKET=%t", !exists),
		fmt.Sprintf("AGC_USE_PUBLIC_SUBNETS=%t", o.publicSubnets),
		fmt.Sprintf("AGC_VERSION=%s", version.Version),
	}
	if o.vpcId != "" {
		environmentVars = append(environmentVars, fmt.Sprintf("VPC_ID=%s", o.vpcId))
	}

	homeDir, err := osutils.DetermineHomeDir()
	if err != nil {
		return err
	}

	cdkAppPath := filepath.Join(homeDir, cdkCoreDir)
	err = o.cdkBootstrap(cdkAppPath, environmentVars)
	if err != nil {
		return err
	}

	return o.deployCoreInfrastructure(cdkAppPath, environmentVars)
}

func (o accountActivateOpts) generateDefaultBucket() (string, error) {
	account, err := o.stsClient.GetAccount()
	if err != nil {
		return "", err
	}
	return generateBucketName(account, o.region), nil
}

func (o accountActivateOpts) cdkBootstrap(cdkAppPath string, environmentVars []string) error {
	progressStream, err := o.cdkClient.Bootstrap(cdkAppPath, environmentVars, bootstrapKey)
	if err != nil {
		return err
	}
	return displayProgress(progressStream, "Bootstrapping CDK...")
}

func (o accountActivateOpts) deployCoreInfrastructure(cdkAppPath string, environmentVars []string) error {
	progressStream, err := o.cdkClient.DeployApp(cdkAppPath, environmentVars, activateKey)
	if err != nil {
		return err
	}
	return displayProgress(progressStream, "Activating account...")
}

func displayProgress(progressStream cdk.ProgressStream, displayMsg string) error {
	if logging.Verbose {
		var lastEvent cdk.ProgressEvent
		for event := range progressStream {
			if event.Err != nil {
				for _, line := range lastEvent.Outputs {
					log.Error().Msg(line)
				}
				return event.Err
			}
			lastEvent = event
		}
	} else {
		return progressStream.DisplayProgress(displayMsg)
	}
	return nil
}

// BuildAccountActivateCommand builds the command for activating AGC in an AWS account.
func BuildAccountActivateCommand() *cobra.Command {
	vars := accountActivateVars{}
	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate AGC in an AWS account.",
		Long: `Activate AGC in an AWS account.
AGC will use your default AWS credentials to deploy all AWS resources
it needs to that account and region.`,
		Example: `
Activate AGC in your AWS account with a custom S3 bucket and VPC.
/code $ agc account activate --bucket my-custom-bucket --vpc my-vpc-id`,
		Args: cobra.NoArgs,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			if vars.publicSubnets && vars.vpcId != "" {
				return clierror.New("account activate", vars, &actionableerror.Error{
					Cause:           fmt.Errorf("both %[1]q and %[2]q cannot be specified together, as %[2]q involves creating a minimal VPC", accountVpcFlag, publicSubnetsFlag),
					SuggestedAction: "Remove one or both of these flags",
				})
			}
			opts, err := newAccountActivateOpts(vars)
			if err != nil {
				return err
			}
			log.Info().Msgf("Activating AGC with bucket '%s' and VPC '%s'", opts.bucketName, opts.vpcId)
			if err := opts.Execute(); err != nil {
				return clierror.New("account activate", vars, err)
			}
			return nil
		}),
	}
	cmd.Flags().StringVar(&vars.bucketName, accountBucketFlag, "", accountBucketFlagDescription)
	cmd.Flags().StringVar(&vars.vpcId, accountVpcFlag, "", accountVpcFlagDescription)
	cmd.Flags().BoolVar(&vars.publicSubnets, publicSubnetsFlag, false, publicSubnetsFlagDescription)
	return cmd
}
