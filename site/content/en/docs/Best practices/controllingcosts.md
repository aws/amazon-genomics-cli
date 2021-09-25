---
title: "Controlling Costs"
date: 2021-09-17T18:01:56-04:00
draft: false
description: >
  Monitoring costs and design considerations to reduce costs
---

When you begin to run large scale workflows frequently it will become important to be able to understand the costs involved and
how to optimize your workflow and use of Amazon Genomics CLI to reduce costs.

## Use AWS Cost Explorer to Report on Costs

AWS Cost Explorer has an easy-to-use interface that lets you visualize, understand, and manage your AWS costs and usage over time.
We recommend you use this tool to gain sight into the costs of running your genomics workflows. At the time of writing AWS Cost Explorer
can only be enabled using the AWS Console so Amazon Genomics CLI won't be able to set this up for you. As a first step you will need to [enable cost explorer](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/ce-getting-started.html) for your 
AWS account.

Amazon Genomics CLI will tag the infrastructure it creates with tags. Application, user, project and context tags are all generated as
appropriate and these can be used as [cost allocation tags](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/cost-alloc-tags.html) 
to determine which account costs are coming from Amazon Genomics CLI and which user, context and project.

Within Cost Explorer the Amazon Genomics CLI tags will be referred to as ["User Defined Cost Allocation Tags"](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/custom-tags.html).
Before a tag can be used in a cost report it must be [activated](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/activating-tags.html). Costs associated with
tags are only available for infrastructure used *after* activation of a tag, so it will not be possible to retrospectively
examine costs.


## Optimizing Requested Container Resources

Tasks in a workflow typically run in Docker containers. Depending on the workflow language there will be some kind of `runtime` definition that specifies the
number of vCPUs and amount of RAM allocated to the task. For example, in WDL you could specify

```
  runtime {
    docker: "biocontainers/plink1.9:v1.90b6.6-181012-1-deb_cv1"
    memory: "8 GB"
    cpu: 2
  }
```

The amount of resource allocated to each container ultimately impacts the cost to run a workflow. Optimally allocating
resources leads to cost efficiency.

## Optimize the longest running, and most parallel tasks first

When optimizing a workflow, focus on those tasks that run the longest as well as those
that have the largest number of parallel tasks as they will make up the majority of the workflow runtime and contribute
most to the cost.


## Consider CPU and memory ratios

EC2 workers for Cromwell AWS Batch compute environments are `c`, `m`, and `r` instance families that
have vCPU to memory ratios of 1:2, 1:4 and 1:8 respectively. Engines that run container based workflows will typically attempt to fit containers to instances in
the most optimal way depending on cost and size requirements, or they will delegate this to a service like AWS Batch. Given that a task requiring 16GB of RAM that could make
use of all available CPUs, then to optimally pack the containers you should specify either 2, 4, or 8 vCPU. Other
values could lead to inefficient packing meaning the resources of the EC2 container instance will be paid for but
not optimally used.

>NOTE: Fully packing an instance can result in it becoming unresponsive if the tasks in the containers use 100%
(or more if they start swapping) of the allocated resources. The instance may then be unresponsive to its management services or the workflow engine and may
time out. To avoid this, always allow for a little overhead, especially in the smaller instances.

The largest instance types deployed by default are from the `4xlarge` size which have 16 vCPU and up to 128 MB of RAM.

## Consider splitting tasks that pipe output

If a workflow task consists of a process that pipes `STDOUT` to another process then both processes will run in the same
container and receive the same resources. If one task requires more resources than the other this might be inefficient, 
it may be better divided into two tasks each with its own `runtime` configuration. Note that this will require the
intermediate `STDOUT` to be written to a file and copied between containers so if this output is very large then keeping
the processes in the same task may be more efficient. Piping very large outputs may a lot of memory so
your container will need an appropriate allocation of memory.

## Use the most cost-effective instance generation

When you specify the `instanceTypes` in a context, as opposed to letting Amazon Genomics CLI do it for you, consider the cost and performance of the instance types with respect to your workflow requirements.
Fifth generation EC2 types (`c5`, `m5`, `r5`) have a lower on-demand price and have higher clock speeds than their 4th
generation counterparts (`c4`, `m4`, `r4`). Therefore, for on-demand compute environments, those instance types should be
preferred. In spot compute environments we suggest using both 4th and 5th generation types as this increases the pool of
available types meaning Batch will be able to choose the instance type that is cheapest and least likely to be
interrupted.

## Deploy Amazon Genomics CLI where your S3 data is

Genomics workflows may need to access considerable amounts of data stored in S3. Although S3 uses global namespaces, buckets
do reside in regions. If you access a lot of S3 data it makes sense to deploy your Amazon Genomics CLI infrastructure in the same region
to avoid cross region data charges.

Further, if you use a custom VPC we recommend deploying a VPC endpoint for S3 so that you do no incur NAT Gateway charges
for data coming from the same region. If you do not you might find that NAT Gateway charges are the largest part of your
workflow run costs. If you allow Amazon Genomics CLI to create your VPC (the default), appropriate VPC endpoints will be setup for you.
Note that VPC endpoints cannot avoid cross region data charges, so you will still want to deploy in the region where most of
your data resides.

## Use Spot Instances

The use of Spot instances can significantly reduce costs of running workflows. However, spot instances may be interrupted when EC2 demand
is high. Some workflow engines, such as Cromwell, can support retries of tasks that fail due to Spot interruption (among other things).
To enable this for Cromwell, include the `awsBatchRetryAttempts` parameter in the `runtime` section of a WDL task with an 
integer number of attempts. 

Even with retries, there is a risk that spot interruption will case a task or entire workflow to fail. Use of an engines call caching capabilities (if available)
can help avoid repeating work if a partially complete workflow needs to be restarted du to Spot instance interruption.

## Use private ECR registries

Each task in a workflow requires access to a container image, and some of these images can be several GB if they contain
large packages like GATK. This can lead to large NAT Gateway traffic charges. To avoid these charges, we recommend deploying
copies of frequently used container images into your accounts private ECR registry.

Amazon Genomics CLI deployed VPCs use a VPC gateway to talk to private ECR registries in your account thereby avoiding NAT Gateway traffic. The gateway is
limited to registries in the same region as the VPC, so to avoid cross-region traffic you should deploy images into the region(s) that you
use for Amazon Genomics CLI.
