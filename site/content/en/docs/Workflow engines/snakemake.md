---
title: "Snakemake"
date: 2021-10-01T17:27:31-04:00
draft: false
weight: 30
description: >
    Details on the Snakemake engine deployed by Amazon Genomics CLI
---

## Description

[Snakemake](https://snakemake.readthedocs.io/en/stable/) is free open source software distributed under the MIT licence 
developed by [Johannes Köster and their team](https://snakemake.readthedocs.io/en/stable/project_info/authors.html). 

The source code for snakemake is available on [GitHub](https://github.com/snakemake/snakemake). When deployed with
Amazon Genomics CLI snakemake uses Batch to distribute the underlying tasks.

## Architecture

There are four components of a snakemake engine as deployed in an Amazon Genomics CLI context:

### WES Adapter

Amazon Genomics CLI communicates with the snakemake engine via a GA4GH [WES](https://github.com/ga4gh/workflow-execution-service-schemas) REST service. The WES Adapter implements
the WES standard and translates WES calls into calls to the snakemake head process.

### Engine Batch Job

For every workflow submitted, the WES adapter will create a new AWS Batch Job that contains the snakemake process responsible
for running that workflow. These snakemake "head" jobs are run in an "On-demand" AWS Fargate compute environment even when the actual workflow
tasks run in a Spot environment. This is to prevent Spot interruptions from terminating the workflow coordinator. 

### Compute Environment

Workflow tasks are submitted by the snakemake head job to an AWS Batch queue and run in containers using an AWS Compute Environment.
Container characteristics are defined by the resources requested in the workflow configuration. AWS Batch coordinates the elastic provisioning of EC2 instances (container hosts)
based on the available work in the queue. Batch will place containers on container hosts as space allows.

#### EFS scratch space and S3 localization

Any context with a snakemake engine will use an Amazon Elastic File System (EFS) volume as scratch space. Inputs from the workflow
are localized to the volume by jobs that the snakemake engine spawns to copy these files to the volume. Outputs are copied back 
to S3 after the workflow is complete. Workflow tasks access the EFS volume to obtain inputs and write intermediates and outputs.

The EFS volume can used by all snakemake engine "head" jobs to store metadata necessary for dependency caching by specifying an argument 
for the conda workspace that is common across all executions. An example of this is `--conda-prefix /mnt/efs/snakemake/conda`.

The EFS volume will remain in your account for the lifetime of the context and are destroyed when contexts are destroyed.
Because the volume will grow in size as you run more workflows we recommend destroying the context when done to avoid on going EFS
charges.

## Using Snakemake as a Context Engine

You may declare snakemake to be the `engine` for any contexts `snakemake` type engine. For example:

```yaml
contexts:
  onDemandCtx:
    requestSpotInstances: false
    engines:
      - type: snakemake
        engine: snakemake
```

## Conda Dependency Caching

Dependency caching is disabled by default so that each workflow can be run independently. If you would like workflow
runs to re-use the Conda cache then please specify a folder under "/mnt/efs" which is where the EFS storage space is
attached. This will enable snakemake to re-use the dependency which will decrease the time that subsequent workflow runs
will take.

To disable call caching you can provide the `--conda-prefix` engine option. You may do this in a workflows `MANIFEST.json` by
adding the following key/ value pair.

```
  "engineOptions": "-j 10 --conda-prefix /mnt/efs/snakemake/conda"
```
