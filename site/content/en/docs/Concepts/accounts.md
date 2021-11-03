---
title: "Accounts"
date: 2021-09-02T13:52:05-04:00
draft: false
weight: 1
description: >
  How AWS Genomics CLI interacts with AWS Accounts
---
Amazon Genomics CLI requires an AWS account in which to deploy the cloud infrastructure required to run and manage workflows. To begin
working with Amazon Genomics CLI and account must be "Activated" by the Amazon Genomics CLI application using the [account activate]( {{< relref "#activate" >}}) command.

## Which AWS Account is Used by Amazon Genomics CLI?

Amazon Genomics CLI uses the same [AWS credential chain](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-precedence) 
used by the AWS CLI to determine what account should be used and with what credentials.
All that is required is that you have an existing AWS account (or create a new one) which contains at least one IAM Principal 
(User/ Role) that you have can access.

## Which Region is Used by Amazon Genomics CLI?

Much like accounts and credentials, Amazon Genomics CLI uses the same chain used by the AWS CLI to determine the region that is being targeted.
For example, if your AWS profile uses `us-east-1` then Amazon Genomics CLI will use the same. Likewise, if you set the `AWS_REGION` environment
variable to `eu-west-1` then that region will be used by Amazon Genomics CLI for all subsequent commands in that shell.

## Shared Infrastructure

When a region is first activated for Amazon Genomics CLI, some basic infrastructure is deployed including a [VPC](https://docs.aws.amazon.com/vpc/latest/userguide/index.html) 
and an [S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/index.html) bucket. This
core infrastructure will be shared by all Amazon Genomics CLI users and projects in that region. Note that context specific infrastructure
is not shared and is unique and namespaced by user and project.

## Bring your Own VPC and S3 Bucket

During account [activation]( {{< relref "#activate" >}}) you may specify an existing VPC ID or S3 bucket name for use by Amazon Genomics CLI. If you do not these will 
be created for you. Although we use AWS best practices for these, if your organization has specific security requirements 
for networking and storage this may be the easiest way to activate Amazon Genomics CLI in your environment.

## Account Commands

A full reference of the account commands is [here]( {{< relref "../Reference/agc_account" >}} )

### `activate`

You can activate an account using `agc account activate`. An account must be activated before any contexts can be deployed
or workflows run. 

Amazon Genomics CLI requires an S3 bucket to store workflow results and associated information. If you prefer to use an existing bucket
you can use the form  `agc account activate --bucket my-existing-bucket`. If you do this the AWS [IAM](https://docs.aws.amazon.com/IAM/latest/UserGuide/index.html) role used to run
Amazon Genomics CLI must be able to write to that bucket.

To use an existing VPC you can use the form  `agc account activate --vpc my-existing-vpc-id`. This VPC must have at least
3 availability zones each with at least one private subnet. The private subnets must have connectivity to the internet, 
such as via a NAT gateway, and connectivity to AWS services either through VPC endpoints or the internet. Amazon Genomics CLI will not
modify the network topology of the specified VPC.

Issuing account activate commands more than once effectively updates the core infrastructure with the difference between
the two commands. For example, if you had previously activated the account using `agc account activate` and later invoked
`agc account activate --bucket my-existing-bucket --vpc my-existing-vpc-id` then Amazon Genomics CLI will update to use `my-existing-bucket`
and the identified VPC. The old VPC and related infrastructure will be destroyed. S3 buckets will be *retained* according
to their retention policy.

If you initially activated the account with `agc account activate --bucket my-existing-bucket --vpc my-existing-vpc-id`
and later invoked `agc account activate` then Amazon Genomics CLI will stop using the previous specified bucket and VPC. *ALL* of the 
pre-existing S3 and VPC infrastructure will be retained and a new bucket and VPC will be created for use by Amazon Genomics CLI.

### `deactivate`

The `deactivate` command is used to remove the core infrastructure deployed by Amazon Genomics CLI in the current region when an 
account is activated. The S3 bucket deployed by Amazon Genomics CLI and its contents are retained. If a VPC and/ or S3 bucket were 
specified by the user during account activation these will also be retained. Any CloudWatch logs produced by Amazon Genomics CLI will
also be retained.

If there are existing deployed contexts the command will fail, however, you can force the removal of these at the same
time with the `--force` flag. Note that this will also interrupt any running workflow of any user in that region.

The deactivate command will only operate on infrastructure in the current region.

If the deployed infrastructure has been modified through the console or the AWS CLI rather than through Amazon Genomics CLI deactivation
may fail due to the infrastructure state being inconsistent with the [CloudFormation](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/index.html) state. If this happens you may need
to manually clean up through the CloudFormation console.

## Costs

Core infrastructure deployed for Amazon Genomics CLI is [tagged]( {{< relref "namespaces#tags" >}} ) with the `application-name: agc` tag. This tag can be activated for cost
tracking in [AWS CostExplorer](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/ce-what-is.html). The core infrastructure is shared and *not* tagged with any username, context name or 
project name.

While an account region is activated there will be ongoing charges from the core infrastructure deployed including things such 
as VPC NAT gateways and VPC Endpoints. If you no longer use Amazon Genomics CLI in a region we recommend you deactivate it. You may also
wish to remove the S3 bucket along with its objects as well as the [CloudWatch](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/index.html) logs produced by Amazon Genomics CLI. These are retained
by default so that you can view workflow results and logs even after deactivation.

### Network traffic

When running genomics workflows, network traffic can become a significant expense when the traffic is routed
through NAT gateways into private subnets (where worker nodes are usually located). To minimize these costs we recommend
the use of VPC Enpoints [(see below)]( {{< relref "#VPC Endpoints" >}} ) as well as activating Amazon Genomics CLI and running your workflows in the same region as your S3
bucket holding your genome files. VPC Gateway endpoints are regional so cross region S3 access will *not* be routed through
a VPC gateway.

If you make use of large container images in your workflows (such as the GATK images) we recommend copying these to a 
private [ECR](https://docs.aws.amazon.com/AmazonECR/latest/userguide/index.html) repository in the same region that you will run your analysis to use ECR endpoints and avoid traffic through
NAT gateways.

### VPC Endpoints

When Amazon Genomics CLI creates a VPC it creates the following VPC endpoints:

* `com.amazonaws.{region}.ecr.api`
* `com.amazonaws.{region}.ecr.dkr`
* `com.amazonaws.{region}.ecs`
* `com.amazonaws.{region}.ecs-agent`
* `com.amazonaws.{region}.ecs-telemetry`
* `com.amazonaws.{region}.logs`
* `com.amazonaws.{region}.s3`

If you provide your own VPC we recommend that the VPC also has these endpoints. This will improve the security posture of
Amazon Genomics CLI in your VPC and will also reduce NAT gateway traffic charges which can be substantial for genomics analyses that use
large S3 objects and/ or large container images.

## Technical Details.

Amazon Genomics CLI core infrastructure is defined in code and deployed by [AWS CDK](https://aws.amazon.com/cdk/). The CDK app responsible for creating the core
infrastructure can be found in [`packages/cdk/apps/core/`](https://github.com/aws/amazon-genomics-cli/tree/main/packages/cdk/apps/core).

