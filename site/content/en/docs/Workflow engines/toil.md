---
title: "Toil"
date: 2022-04-26T15:34:00-04:00
draft: false
weight: 20
description: >
    Details on the Toil engine (CWL mode) deployed by Amazon Genomics CLI
---

## Description

[Toil](http://toil.ucsc-cgl.org/) is a workflow engine developed by the
[Computational Genomics Lab](https://cglgenomics.ucsc.edu/) at the
[UC Santa Cruz Genomics Institute](https://genomics.ucsc.edu/). In Amazon Genomics
CLI, Toil is an engine that can be deployed in a
[context]( {{< relref "../Concepts/contexts" >}} ) as an
[engine]( {{< relref "../Concepts/engines">}} ) to run workflows written in the
[Common Workflow Language](https://www.commonwl.org/) (CWL) standard, version
[v1.0](https://www.commonwl.org/v1.0/), [v1.1](https://www.commonwl.org/v1.1/),
and [v1.2](https://www.commonwl.org/v1.2/) (or mixed versions).

Toil is an open source project distributed by UC Santa Cruz under the [Apache 2
license](https://github.com/DataBiosphere/toil/blob/master/LICENSE) and
available on
[GitHub](https://github.com/DataBiosphere/toil).

## Architecture

There are two components of a Toil engine as deployed in an Amazon Genomics
CLI context:

### Engine Service

The Toil engine is run in "server mode" as a container service in ECS. The
engine can run multiple workflows asynchronously. Workflow tasks are run in an
elastic [compute environment]( #compute-environment ) and monitored by Toil.
Amazon Genomics CLI communicates with the Toil engine via a GA4GH
[WES](https://github.com/ga4gh/workflow-execution-service-schemas) REST service
which the server offers, available via API Gateway.

### Compute Environment

Workflow tasks are submitted by Toil to an AWS Batch queue and run in
Toil-provided containers using an AWS Compute Environment. Tasks which use the
[CWL `DockerRequirement`](https://www.commonwl.org/user_guide/07-containers/index.html)
will additionally be run in sibling containers on the host Docker daemon. AWS
Batch coordinates the elastic provisioning of EC2 instances (container hosts)
based on the available work in the queue. Batch will place containers on
container hosts as space allows.

#### Disk Expansion

Container hosts in the Batch compute environment use EBS volumes as local
scratch space. As an EBS volume approaches a capacity threshold, new EBS
volumes will be attached and merged into the file system. These volumes are
destroyed when AWS Batch terminates the container host. CWL disk space
requirements are ignored by Toil when running against AWS Batch.

This setup means that workflows that succeed on AGC may fail on other CWL
runners (because they do not request enough disk space) and workflows that
succeed on other CWL runners may fail on AGC (because they allocate disk space
faster than the expansion process can react).


