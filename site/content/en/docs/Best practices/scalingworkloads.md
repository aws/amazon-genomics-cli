---
title: "Scaling Workloads"
date: 2021-09-17T17:59:52-04:00
draft: false
description: >
    Making workflows run at scale
---

Workflows with considerable compute requirements can incur large costs and may fail due to infrastructure constraints.
The following considerations will help you design workflows that will perform better at scale.

## Large compute requirements

By default, contexts created by AGC will allocate compute nodes with a size of up to `4xlarge`. These types have 16 vCPU
and up to 128 GB of RAM. If an individual task requires additional resources you may specify these
in the `instanceTypes` array of the project context. For example:

```yaml
contexts:
    prod:
        requestSpotInstances: false
        instanceTypes:
            - c5.16xlarge
            - r5.24xlarge
```

## Large data growth

When using the Nextflow or Cromwell engines the EC2 container instances that carry out the work use a script to detect
and automatically expand disk capacity. Generally, this will allow disk space to increase to the amount required to hold
inputs, scratch and outputs. However, it can take up to a minute to attach new storage so events that fill disk space
in under a minute can result in failure.

## Large numbers of inputs/ outputs

Typically, genomics files are large and best stored in S3. However, most applications used in genomics workflows cannot
read directly from S3. Therefore, these inputs must be localized from S3. Compute work will not be able to begin until localization is complete
so "divide and conquer" strategies are useful in these cases. 

Whenever possible compress inputs (and outputs) appropriately. The CPU overhead of compression will be low compared to the
network overhead of localization and delocalization.

Localization of large numbers of large files from S3 will put load on the network interface of the worker nodes and the
node may experience transient network failures or S3 throttling. While we have included retry-with-backoff logic for
localization it is not impossible that downloads may occasionally fail. Failures (and retries) will be recorded in 
the workflow task logs.

## Parallel Steps

Workflows often contain parallel steps where many individual tasks are computed in parallel. AGC makes use of elastic compute
clusters to scale to these requirements. Each context will deploy an elastic compute cluster with a minimum of 0 vCPU and a maximum of 256 vCPU. No individual task
may use more than 256 vCPU. Smaller tasks may be run in parallel up to the maximum of 256 vCPU. Once that limit is met, additional
tasks will be queued to run when capacity becomes free.

Each parallel task is isolated meaning each task will need a local copy of its inputs. When large numbers of parallel tasks,
require the same inputs (for example reference genomes) you may observe contention for network resources and transient S3
failures. While we have included retry with backoff logic we recommend keeping the number of parallel tasks requiring the same inputs below 500. Fewer,
if the tasks inputs are large.

An extreme example is Joint Genotyping. This type of analysis benefits from processing large numbers of samples at the same.
Further, the user may wish to genotype many intervals concurrently. Finally, the step of merging the variant calls will 
import the variants from all intervals. In our experience, a naive implementation calling 100 samples over 100 intervals
is feasible. Also, feasible is calling ~20 samples over 500 intervals. At larger scales it would be worth considering dividing
tasks by chromosome or batching inputs.

## Container throttling

Some container registries will throttle container access from anonymous accounts. Because each task in a workflow uses
a container large or frequently run workflows may not be able to access their required containers. While compute clusters
deployed by AGC are configured to cache containers this is only available on a per-instance basis. Further, due to the 
elastic nature of the clusters instances with cached container images are frequently shutdown. All of this will potentially
lead to an excess of requests. To avoid this we recommend using registries that don't impose these limits, or using images
hosted in an ECR registry in your AWS account.
