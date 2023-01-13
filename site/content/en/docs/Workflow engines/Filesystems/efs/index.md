---
title: "EFS Workflow Filesystem"
date: 2022-03-11T15:21:15-05:00
draft: False
---

## Amazon EFS Workflow Filesystem

Workflow engines that support it may use Amazon EFS as a shared "scratch" space for hosting workflow intermediates and 
outputs. Initial inputs are localized once from S3 and final outputs are written back to S3 when the workflow is complete.
All intermediate I/O is performed against the EFS filesystem.

### Advantages

1. Compared with the [S3 Filesystem]( {{< relref "../s3" >}} ) there is no redundant I/O of inputs from S3.
2. Each tasks individual I/O operations tend to be smaller than the copy from S3 so there is less network congestion on the container host.
3. Option to use provisioned IOPs to provide high sustained throughput.
4. The volume is elastic and will expand and contract as needed.
5. It is simple to start an Amazon EC2 instance from the AWS console and connect it to the EFS volume to view outputs as they are created. This can be useful for debugging a workflow.

### Disadvantages

1. Amazon EFS volumes are more expensive than storing intermediates and output in S3, especially when the volume uses provisioned throughput.
2. The volume exists for the lifetime of the context and will incur costs based on its size for the lifetime of the context. If you no longer need the context we recommend destroying it.
3. Call caching is only possible for as long as the volume exists, i.e. the lifetime of the context.

### Provisioned Throughput

Amazon EFS volumes deployed by the Amazon Genomics CLI use ["bursting"](https://docs.aws.amazon.com/efs/latest/ug/performance.html#bursting) 
throughput by default. For workflows that have high I/O throughput or in scenarios where you may have many workflows 
running in the same context at the same time, you may exhaust the burst credits of the volume. 
This might cause a workflow to slow down or even fail. Available volume credits can be [monitored](https://docs.aws.amazon.com/efs/latest/ug/monitoring_overview.html)
in the Amazon EFS console, and/ or Amazon CloudWatch. 

If you observe the exhaustion of burst credits you may want to consider
deploying a context with [provisioned throughput](https://docs.aws.amazon.com/efs/latest/ug/performance.html#provisioned-throughput).
Throughput is provisioned in MiB/s and may be upto 1024 MiB/s. Note that this is an **additional expense** for the EFS volume
and is charged even when the volume has no data stored in it. If you choose an EFS volume with provisioned throughput we encourage
you to destroy the context whenever it is not in use to minimize costs. To determine the costs of provisioned throughput
you may use the [AWS Price Calculator for EFS](https://calculator.aws/#/addService/EFS)

The following fragment of an `agc-project.yaml` file is an example of how to configure provisioned throughput for the
Amazon EFS volume used by miniwdl in an Amazon Genomics CLI context.

```yaml
myContext:
    engines:
      - type: wdl
        engine: miniwdl
        filesystem:
          fsType: EFS
          configuration:
            provisionedThroughput: 100
```

### Supporting Engines

The use of Amazon EFS as a shared file system is supported by the [miniwdl]( {{< relref "../../miniwdl" >}} ) and 
[Snakemake]( {{< relref "../../snakemake" >}} ) engines. Both use EFS with bursting throughput by default and both 
support provisioned throughput.