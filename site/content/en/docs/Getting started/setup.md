---
title: "Setup"
date: 2021-09-07T13:42:50-04:00
draft: false
weight: 30
---


## Account activation

To start using Amazon Genomics CLI with your AWS account, you need to activate it.

```
agc account activate
```

This will create the core infrastructure that Amazon Genomics CLI needs to operate, which includes a DynamoDB table, an S3 bucket and a VPC. This will take ~5min to complete. You only need to do this once per account region.

The DynamoDB table is used by the CLI for persistent state. The S3 bucket is used for durable workflow data and AGC metadata and the VPC is used to isolate compute resources. You can specify your own preexisting S3 Bucket or VPC if needed using `--bucket` and `--vpc` options.

## Define a username

AGC requires that you define a username and email. You can do this using the following command:

`agc configure email you@youremail.com`

The username only needs to be configured once per computer that you use AGC from.
