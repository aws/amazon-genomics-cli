---
title: "miniwdl"
date: 2021-10-01T17:27:31-04:00
draft: false
weight: 30
description: >
    Details on the miniwdl engine deployed by Amazon Genomics CLI
---

## Description

[miniwdl](https://miniwdl.readthedocs.io/en/latest/index.html) is free open source software distributed under the MIT licence 
developed by the [Chan Zuckerberg Initiative](https://chanzuckerberg.com/). 

The source code for miniwdl is available on [GitHub](https://github.com/chanzuckerberg/miniwdl). When deployed with
Amazon Genomics CLI miniwdl makes use of the [miniwdl-aws extension](https://github.com/miniwdl-ext/miniwdl-aws) which is
also distributed under the MIT licence.

## Architecture

There are four components of a miniwdl engine as deployed in an Amazon Genomics CLI context:

### WES Adapter

Amazon Genomics CLI communicates with the miniwdl engine via a GA4GH [WES](https://github.com/ga4gh/workflow-execution-service-schemas) REST service. The WES Adapter implements
the WES standard and translates WES calls into calls to the miniwdl head process.

### Engine Batch Job

For every workflow submitted, the WES adapter will create a new AWS Batch Job that contains the miniwdl process responsible
for running that workflow. These miniwdl "head" jobs are run in an "On-demand" AWS Fargate compute environment even when the actual workflow
tasks run in a Spot environment. This is to prevent Spot interruptions from terminating the workflow coordinator. 

### Compute Environment

Workflow tasks are submitted by the miniwdl head job to an AWS Batch queue and run in containers using an AWS Compute Environment.
Container characteristics are defined by the resources requested in the workflow configuration. AWS Batch coordinates the elastic provisioning of EC2 instances (container hosts)
based on the available work in the queue. Batch will place containers on container hosts as space allows.

#### EFS scratch space and S3 localization

Any context with a miniwdl engine will use an Amazon Elastic File System (EFS) volume as scratch space. Inputs from S3 are
localized to the volume by jobs that the miniwdl engine spawns to copy these files to the volume. Outputs are copied back 
to S3 using a similar process. Workflow tasks access the EFS volume to obtain inputs and write intermediates and outputs.

The EFS volume is used by all miniwdl engine "head" jobs to store metadata necessary for call caching.

The EFS volume will remain in your account for the lifetime of the context and are destroyed when contexts are destroyed.
Because the volume will grow in size as you run more workflows we recommend destroying the context when done to avoid on going EFS
charges.

## Using miniwdl as a Context Engine

You may declare miniwdl to be the `engine` for any contexts `wdl` type engine. For example:

```yaml
contexts:
  onDemandCtx:
    requestSpotInstances: false
    engines:
      - type: wdl
        engine: miniwdl
```

## Call Caching

Call caching is enabled by default for miniwdl and because the metadata is stored in the contexts EFS volume call caching
will work across different engine "head" jobs.

To disable call caching you can provide the `--no-cache` engine option. You may do this in a workflows `MANIFEST.json` by
adding the following key/ value pair.

```
  "engineOptions": "--no-cache"
```
