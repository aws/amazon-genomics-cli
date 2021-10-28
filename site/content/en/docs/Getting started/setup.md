---
title: "Setup"
date: 2021-09-07T13:42:50-04:00
draft: false
weight: 30
---

## CDK Bootstrap

Amazon Genomics CLI uses AWS CDK to deploy infrastructure. If you have not already done so, you will need to "bootstrap" CDK into the account and region where you plan to use Amazon Genomics CLI.
CDK Bootstrap deploys the infrastructure needed to allow CDK to deploy CDK defined infrastructure. Full details are available
[here](https://docs.aws.amazon.com/cdk/latest/guide/bootstrapping.html).

Amazon Genomics CLI will install its own local copy of of the AWS CDK. If you do not have the AWS CDK installed globally on your machine use can also use this local copy by changing to the Amazon Genomics CLI CDK installation location:

```
cd $HOME/.agc/cdk
```

and using `npx` to invoke the CDK:

```
npx cdk ...
```

To bootstrap a single region:

`cdk bootstrap aws://ACCOUNT-NUMBER-1/REGION-1`

If you plan to use Amazon Genomics CLI in multiple regions you can use this form of the command:

`cdk bootstrap aws://ACCOUNT-NUMBER-1/REGION-1 aws://ACCOUNT-NUMBER-1/REGION-2 ...`

### Example
```shell
cdk bootstrap 123456789012/us-east-1 123456789012/us-west-1
```


## Account activation

To start using Amazon Genomics CLI with your AWS account, you need to activate it.

```
agc account activate
```

This will create the core infrastructure that Amazon Genomics CLI needs to operate, which includes a DynamoDB table, an S3 bucket and a VPC. This will take ~5min to complete. You only need to do this once per account region.

The DynamoDB table is used by the CLI for persistent state. The S3 bucket is used for durable workflow data and Amazon Genomics CLI metadata and the VPC is used to isolate compute resources. You can specify your own preexisting S3 Bucket or VPC if needed using `--bucket` and `--vpc` options.

## Define a username

Amazon Genomics CLI requires that you define a username and email. You can do this using the following command:

`agc configure email you@youremail.com`

The username only needs to be configured once per computer that you use Amazon Genomics CLI from.
