---
title: "Users"
date: 2021-09-03T09:59:23-04:00
draft: false
weight: 5
description: >
  How Amazon Genomics CLI identifies users
---

When the CLI is set up, the user of the CLI is defined using the `agc configure email` command. This email should be 
unique to the individual user. This email address is used to determine a unique user ID which will be used to uniquely
identify infrastructure belonging to that user.

## Amazon Genomics CLI Users are Not IAM Users (or Principals) 
Amazon Genomics CLI users are primarily used for identification and as a component of namespacing. They are not a security measure, nor 
are they related to IAM users or roles. All AWS activities carried out by Amazon Genomics CLI will be done using the AWS credentials in
the environment where the CLI is installed and are *not* based on the Amazon Genomics CLI username.

For example. If Amazon Genomics CLI is installed on an EC2 instance and configured with the email `someone@company.com` Amazon Genomics CLI will interact
with AWS resources based solely on the IAM Role assigned to that EC2 via it's instance profile. Like wise if you use Amazon Genomics CLI
on your laptop then the IAM role that you use will be determined by the same process as is used by the AWS CLI.

## Who am I?

To find out what username and email has been configured in your current environment you can use the following command:
```shell
agc configure describe
```

## Changing user

If you update your configured email, a new user identity is generated. If this is done while infrastructure is deployed,
Amazon Genomics CLI may no longer be able to identify that infrastructure as belonging to your project. We strongly 
recommend stopping all running workflows and destroying all your deployed contexts from *all* projects before changing user.
If you do not do this, you or an account administrator will need to identify any un-needed infrastructure in the CloudFormation
console and remove it from there.