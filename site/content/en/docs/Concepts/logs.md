---
title: "Logs"
date: 2021-08-31T17:30:49-04:00
draft: false
weight: 50
description: >
    Logs are produced by contexts, engines and workflow tasks. Understanding how to access them is critical to monitoring and debugging workflows.
---

The infrastructure deployed by Amazon Genomics CLI records logs for many activities including the workflow runs, workflow
engines as well as infrastructure. The logs are recorded in CloudWatch but are accessible through the CLI.

When debugging or reviewing a workflow run, the engine logs and workflow logs will be the most useful. For diagnosing
infrastructure or access problems the adapter logs and access logs will be informative.

## Engine Logs

Engine logs are the logs produced by a workflow engine in a context. The logs produced depend on the engine implementation.
Engines that run in "server" mode, such as Cromwell, produce a single log for the lifetime of the context that encompass
all workflows run through that engine. Engines that run as "head node" will produce discrete engine logs for each run.

## Workflow Logs

Workflow logs are the aggregate logs for all steps in a workflow run (instance). Any workflow steps that are retrieved from
a call cache are not run so there will be no workflow logs for these steps. Consulting the engine logs may show details of
the call cache. If a previously successful workflow is run with no changes in inputs or parameters it may have all steps
retrieved from the cache in which case there will be no workflow logs although the workflow instance will be marked as a 
success and engine logs will be produced. The outputs for a completely cached workflow will also be available.

## Adapter Logs

Adapter logs consist of any logs produced by a WES adapter for a workflow engine. They can reveal information such as
the WES API calls that are made to the engine by AGC and any errors that may have occurred. 

## Access Logs

AGC talks to an engine via
API Gateway which routes to the WES adapter. If an expected call does not appear in the adapter logs it may have been
blocked or incorrectly routed in the API Gateway. The API Gateway access logs may be informative in this case.

## Commands

A full reference of AGC `logs` commands are available [here]( {{< relref "../Reference/agc_logs" >}} )

## Costs

Amazon Genomics CLI logs are stored in CloudWatch and accessed using the CloudWatch APIs. Standard CloudWatch charges apply.
All logs are retained permanently, even after a context is destroyed and AGC removed from an account. If they are no longer needed they may be removed
using the AWS Console or the AWS CLI.