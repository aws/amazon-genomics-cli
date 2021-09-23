---
title: "Data"
date: 2021-08-31T17:27:22-04:00
draft: false
weight: 30
description: >
    Data sets
---
To run an analysis you need data. In the `agc-project.yaml` file of an AGC [project]({{< relref "projects" >}}) `data` is a list of data locations 
which can be used by the [contexts]({{< relref "contexts" >}}) of the project.

In the example data definition below we are declaring that the project's contexts will be allowed to access the three
listed S3 bucket URIs.

```yaml
data:
  - location: s3://gatk-test-data
    readOnly: true
  - location: s3://broad-references
    readOnly: true
  - location: s3://1000genomes-dragen-3.7.6
    readOnly: true
```

The contexts of the project will be *denied* access to all other S3 location except for the S3 bucket created or associated
when the [account]( {{< relref "accounts" >}} ) was initialized by AGC.

Declaring access in the project will only ensure your infrastructure is correctly configured to access the bucket. If
the target bucket is further restricted, such as by an access control list or bucket policy, you will still be denied access.
In these cases you should work with the bucket owner to facilitate access.

### Read and Write

The default value of `readOnly` is `true`. At the time of writing, write access is not supported (except for the AGC core S3 bucket)

### Access to a Prefix

The above examples will grant read access to an entire bucket. You can grant more granular access to a prefix within a bucket,
for example:

```yaml
data:
  - location: s3://my-bucket/my/prefix/
```

### Cross Account Access

A bucket in another AWS account can be accessed if the owner has set up access, and you are using a role that is allowed access.
See [cross account access](https://aws.amazon.com/premiumsupport/knowledge-center/cross-account-access-s3/) for details.

## Updating Data Sources

If data definitions are added to or removed from a project definition the change will *not* be reflected in deployed contexts
until they are updated. This can be done with `agc context deploy --all` for all contexts or by using a context name to update 
only one. See [`context deploy`]( {{< relref "contexts#deploy" >}} ) for details.

{{% alert title="Warning" color="warning" %}}
Removing access to S3 buckets while there are running workflows in a project may cause the workflow to fail if it depends
on access to data in those buckets.
{{% /alert %}}

## Technical Details

When a context is deployed, IAM roles used by the infrastructure of the context will be granted s3 permissions to perform
some S3 read (or read and write) actions on the listed locations. The permissions are defined in CDK code in `/packages/cdk/apps/`.
The CDK code does not modify any data in the buckets or any other bucket policies or configurations.

