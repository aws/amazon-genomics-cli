---
title: "Engines"
date: 2021-08-31T17:28:39-04:00
draft: false
weight: 45
description: >
    Workflow engines parse and manage the tasks in a workflow
---

A workflow engine is defined as part of a [context]( {{< relref "contexts" >}} ). A context is currently limited to one workflow engine. The workflow engine will manage the execution of any [workflows]( {{< relref "workflows" >}} ) submitted
by Amazon Genomics CLI. When the context is deployed, an endpoint will be made available
to Amazon Genomics CLI through which it will submit workflows and workflow commands to the engine according to the WES API specification.

## Supported Engines and Workflow Languages

Currently, Amazon Genomics CLI's officially supported engines can be used to run the following workflows:

| Engine                                                 | Language                                                        | Language Versions                                                                                                       | Run Mode     |
|--------------------------------------------------------|-----------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|--------------|
| [Cromwell](https://cromwell.readthedocs.io/en/stable/) | [WDL](https://openwdl.org)                                      | All versions up to 1.0                                                                                                  | Server       |
| [Nextflow](https://www.nextflow.io)                    | [Nextflow DSL](https://www.nextflow.io/docs/latest/script.html) | Standard and DSL 2                                                                                                      | Head Process |
| [miniwdl](https://miniwdl.readthedocs.io/en/latest/)   | [WDL](https://openwdl.org)                                      | [documented here](https://miniwdl.readthedocs.io/en/latest/runner_reference.html?highlight=errata#wdl-interoperability) | Head Process |

Overtime we plan to add additional engine and language support and provide the ability for third party developers to 
develop engine plugins.

### Run Mode

#### Server

In server mode the engine runs as a long-running process that exists for the lifetime of the context. All workflow instances sent to the context are handled by that server. The server resides on on-demand instances to prevent Spot interruption even if the workflow tasks are run on Spot instances

#### Head Process

Head process engines are run when a workflow is submitted, manage a single workflow and only run for the lifetime of the workflow. If multiple workflows are submitted to a context in parallel then multiple head processes are spawned. The head processes always run on on-demand resources to prevent Spot interruption even if the workflow tasks are run on Spot instances. 


## Engine Definition

An engine is defined within a `context` definition of the [project YAML file]( {{< relref "projects#project-file-structure" >}} ) file as a map. For example, the following snippet
defines a WDL engine of type `cromwell` as part of the context named `onDemandCtx`. There may be one engine defined 
for each supported language.

```yaml
contexts:
  onDemandCtx:
    requestSpotInstances: false
    engines:
      - type: wdl
        engine: cromwell
```

## Commands

There are no commands specific to engines. Engines are [deployed]( {{< relref "contexts#deploy" >}} ) along with contexts by the [`context` commands]( {{< relref "contexts#context-commands" >}} ) and workflows
are run using the [`workflow` commands]( {{< relref "workflows#commands" >}} ).

## Costs

The costs associated with an engine depend on the actual infrastructure required by the engine. In the case of the Cromwell,
the engine runs in "server" mode as an [AWS ECS Fargate](https://docs.aws.amazon.com/AmazonECS/latest/userguide/index.html) container using an 
[Amazon Elastic File System](https://docs.aws.amazon.com/efs/latest/ug/index.html) for metadata storage. The container
will be running for the entire time the context is deployed, even when no workflows are running. To avoid this cost we
recommend destroying the context when it is not needed. The Nextflow engine runs as a single batch job per workflow instance
and is only running when workflows are running.

In both cases a serverless WES API endpoint is deployed through [Amazon API Gateway](https://docs.aws.amazon.com/apigatewayv2/latest/api-reference/) to act as the interface between Amazon Genomics CLI and
the engine. 

### Tags

Being part of a context, engine related infrastructure is [tagged]( {{< relref "namespaces#tags" >}} ) with the context name, username and project name. These tags may be used to help
differentiate costs.

## Technical Details

Supported engines are currently deployed with configurations that allow them to make use of files in S3 and submit workflows
as jobs to AWS Batch. Because the current generation of engines we support do not directly support the [WES API](https://ga4gh.github.io/workflow-execution-service-schemas/docs/), adapters
are deployed as Fargate container tasks. AWS API Gateway is used to provide a gateway between Amazon Genomics CLI and the WES adapters.

When `workflow` commands are issued on Amazon Genomics CLI, it will send WES API calls to the appropriate endpoint. The adapter mapped 
to that endpoint will then translate the WES command and either send the command to the engines REST API for Cromwell, or
spawn a Nextflow engine task and submit the workflow with that task. At this point the engine is responsible for creating
controlling and destroying the resources that will be used for task execution.
