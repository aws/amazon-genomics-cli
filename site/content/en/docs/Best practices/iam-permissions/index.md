---
title: "IAM Permissions"
date: 2022-04-07T09:51:56-04:00
draft: false
description: >
    Minimum IAM Permissions required to use AGC
---

Amazon Genomics CLI is used to deploy and interact with infrastructure in an AWS account. Amazon Genomics CLI will use
the permissions of the current profile to perform its actions. The profile would either be the users profile or, if being 
run from an EC2 instance, the attached profile of the instance. No matter the source of the role it must have sufficient 
permissions to perform the necessary tasks. In addition, best practice recommends that the profile only grant minimal
permissions to maintain security and prevent unintended action.

### Recommended Minimal Permissions

As part of the Amazon Genomics CLI repository we have included a CDK project that can be used by an account administrator 
to generate the necessary minimum policies.

#### Pre-requisites 

Before generating the policies you need to do the following:
1. Install `node` and `npm`. We recommend using node v14.17 installed via `nvm` 
2. Install Amazon CDK (`npm install -g aws-cdk@latest`)
3. An AWS account where you will use Amazon Genomics CLI
4. A role in that account that allows the creation of IAM roles and policies

#### Generate Roles and Policies

1. Clone the Amazon Genomics CLI repository locally: `git clone https://github.com/aws/amazon-genomics-cli.git`
2. cd `amazon-genomics-cli/extras/agc-minimal-permissions/`
3. `npm install`
4. `cdk deploy`

You will see output similar to the following:

```
✨  Synthesis time: 2.91s

AgcPermissionsStack: deploying...
AgcPermissionsStack: creating CloudFormation changeset...

 ✅  AgcPermissionsStack

✨  Deployment time: 44.39s

Stack ARN:
arn:aws:cloudformation:us-east-1:123456789123:stack/AgcPermissionsStack/6ace55f0-b67c-11ec-a5d3-0a1e6da159c9

✨  Total time: 47.3s
```

Using the emitted Stack ARN you can identify the policies created. You can also inspect the stack in the CloudFormation console.

For example:

```shell
aws cloudformation describe-stack-resources --stack-name <stack arn>
```

with output similar to:

```json
{
    "StackResources": [
        {
            "StackName": "AgcPermissionsStack",
            "StackId": "arn:aws:cloudformation:us-east-1:123456789123:stack/AgcPermissionsStack/6ace55f0-b67c-11ec-a5d3-0a1e6da159c9",
            "LogicalResourceId": "CDKMetadata",
            "PhysicalResourceId": "6ace55f0-b67c-11ec-a5d3-0a1e6da159c9",
            "ResourceType": "AWS::CDK::Metadata",
            "Timestamp": "2022-04-07T14:10:30.922000+00:00",
            "ResourceStatus": "CREATE_COMPLETE",
            "DriftInformation": {
                "StackResourceDriftStatus": "NOT_CHECKED"
            }
        },
        {
            "StackName": "AgcPermissionsStack",
            "StackId": "arn:aws:cloudformation:us-east-1:123456789123:stack/AgcPermissionsStack/6ace55f0-b67c-11ec-a5d3-0a1e6da159c9",
            "LogicalResourceId": "agcadminpolicy25003180",
            "PhysicalResourceId": "arn:aws:iam::123456789123:policy/AgcPermissionsStack-agcadminpolicy25003180-1ST0KJ0I5J45R",
            "ResourceType": "AWS::IAM::ManagedPolicy",
            "Timestamp": "2022-04-07T14:10:41.597000+00:00",
            "ResourceStatus": "CREATE_COMPLETE",
            "DriftInformation": {
                "StackResourceDriftStatus": "NOT_CHECKED"
            }
        },
        {
            "StackName": "AgcPermissionsStack",
            "StackId": "arn:aws:cloudformation:us-east-1:123456789123:stack/AgcPermissionsStack/6ace55f0-b67c-11ec-a5d3-0a1e6da159c9",
            "LogicalResourceId": "agcuserpolicy346A2D4F",
            "PhysicalResourceId": "arn:aws:iam::123456789123:policy/AgcPermissionsStack-agcuserpolicy346A2D4F-1X9U4HCQ8Z19U",
            "ResourceType": "AWS::IAM::ManagedPolicy",
            "Timestamp": "2022-04-07T14:10:41.981000+00:00",
            "ResourceStatus": "CREATE_COMPLETE",
            "DriftInformation": {
                "StackResourceDriftStatus": "NOT_CHECKED"
            }
        },
        {
            "StackName": "AgcPermissionsStack",
            "StackId": "arn:aws:cloudformation:us-east-1:123456789123:stack/AgcPermissionsStack/6ace55f0-b67c-11ec-a5d3-0a1e6da159c9",
            "LogicalResourceId": "agcuserpolicycdk27FA61BC",
            "PhysicalResourceId": "arn:aws:iam::123456789123:policy/AgcPermissionsStack-agcuserpolicycdk27FA61BC-OXS49AINGPIG",
            "ResourceType": "AWS::IAM::ManagedPolicy",
            "Timestamp": "2022-04-07T14:10:41.747000+00:00",
            "ResourceStatus": "CREATE_COMPLETE",
            "DriftInformation": {
                "StackResourceDriftStatus": "NOT_CHECKED"
            }
        }
    ]
}
```

Three resources of type `AWS::IAM::ManagedPolicy` are created:

* The resource with a name similar to `agcadminpolicy25003180` identify policies which grant sufficient permission to run `agc account activate` and `agc account deactivate` and should be attached to the profile of users who need to perform that action
* Two resources with names similar to `agcuserpolicy346A2D4F` and `agcuserpolicycdk27FA61BC` identify policies which allow all other Amazon Genomics CLI actions. These should be attached to profiles that will use Amazon Genomics CLI day to day.