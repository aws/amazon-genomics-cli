---
title: "Contexts"
date: 2021-08-31T17:26:49-04:00
draft: false
weight: 20
description: >
    Contexts are the set of cloud resources used to run a workflow
---
## What is a Context?

A context is a set of cloud resources. Amazon Genomics CLI runs [workflows]( {{< relref "workflows" >}} ) in a context. A deployed context will include an 
[engine]( {{< relref "engines" >}}) that can interpret
and manage the running of a workflow along with compute resources that will run the individual tasks of the workflow. The
deployed context will also contain any resources needed by the engine or compute resources including any security, permissions
and [logging]( {{< relref "logs" >}} ) capabilities. Deployed contexts are [namespaced]( {{< relref "namespaces" >}}) based on the user, project and context name so that resources
are isolated, preventing collisions.

When a workflow is run the user will decide which context will run it. For example, you might choose to submit a workflow
to a context that uses "Spot priced" resources or one that uses "On Demand" priced resources.

When deployed context resources that require a VPC will be deployed into the VPC that was specified when the [account]( {{< relref "accounts" >}} ) was
activated.

## How is a Context Defined?

A context is defined in the YAML file that defines the [project]( {{< relref "projects" >}} ). A project has at least one context but may have many.
Contexts must have unique names and are defined as YAML maps.

A context may request use of [Spot priced](https://aws.amazon.com/ec2/spot/pricing/) compute resources with `requestSpotInstances: true`. The default value is `false`.

A context must define an array of one or more `engines`. Each engine definition must specify the workflow language that it 
will interpret. For each language Amazon Genomics CLI has a default engine however, users may specify the exact engine in the `engine`
parameter.

## Context Properties

### Instance Types

You may optionally specify the instance types to be used in a context. This can be a specific type such as `r5.2xlarge`
or it can be an instance family such as `c5` or a combination. By default, a context will use instance types up to `4xlarge`

> Note, if you only specify large instance types you will be using those instances for running even the smallest tasks so
we recommend including smaller types as well.

Ensure that any custom types you list are available in the region that you're using with Amazon Genomics CLI or the 
context will fail to deploy. You can obtain a list using the following command

```shell
aws ec2 describe-instance-type-offerings \
    --region <region_name>
```


#### Examples

The following snippet defines two contexts, one that uses spot resources and one that uses on demand. Both contain a
WDL engine.

```yaml
...
contexts:
  # The on demand context uses on demand EC2 instances which may be more expensive but will not be interrupted
  onDemandCtx:
    requestSpotInstances: false
    engines:
      - type: wdl
        engine: cromwell
```


```yaml
  # The spot context uses EC2 spot instances which are usually cheaper but may be interrupted
  spotCtx:
    requestSpotInstances: true
    engines:
      - type: wdl
        engine: cromwell
...
```

The following context may use any instance type from the `m5`, `c5` or `r5` families

```yaml
contexts:
  nfLargeCtx:
    instanceTypes: [ "c5", "m5", "r5" ]
    engines:
      - type: nextflow
        engine: nextflow
```

### Max vCpus

*default:* 256

You may optionally specify the maximum number of vCpus used in a context. This is the max total amount of vCpus of all the jobs currently 
running within a context. When the max has been reached, additional jobs will be queued.

*note:* if your vCPU limit is lower than maxVCpus then you won't get as many as requested and would need to make a limit increase.
```yaml
contexts:
  largeCtx:
    maxVCpus: 2000
    engines:
      - type: nextflow
        engine: nextflow
```

## Context Commands

A full reference of context commands is [here]( {{< relref "../Reference/agc_context" >}} )

### `describe`

The command `agc context describe <context-name> [flags]` will describe the named context as defined in the project YAML
as well as other relevant account information.

### `list`

The command `agc context list [flags]` will list the names of all contexts defined in the project YAML file

### `deploy`

The command `agc context deploy -c <context-name> [flags]` is used to deploy the cloud infrastructure required by the context.
If the context is already running the existing infrastructure will be updated to reflect changes in project YAML. For example
if you added another `data` definition in your project and run `agc context deploy -c <context-name>` then the deployed context
will be updated to allow access to the new data.

All contexts defined in the project YAML can be deployed or updated using the `--all` flag.

Individually named contexts can be deployed or updated as positional arguments. For example: `agc context deploy -c ctx1 -c ctx2`
will deploy the contexts `ctx1` and `ctx2`.

The inclusion of the `--verbose` flag will show the full CloudFormation output of the context deployment.

### `destroy`

A contexts cloud resources can be "destroyed" using the `agc context destroy -c <context-name>` command. This will remove any 
infrastructure artifacts associated with the context unless they are defined as being retained. Typically, things like logs
and workflow outputs on S3 are retained when a context is destroyed.

All deployed contexts can be destroyed using the `--all` flag.

Multiple contexts can be destroyed in a single command using positional arguments. For example: `agc context destroy -c ctx1 -c ctx2`
will destroy the contexts `ctx1` and `ctx2`.

### `status`

The status command is used to determine the status of a *deployed* context or context instance. This can be useful to determine
if an instance of a particular context is already deployed. It can be used to determine if the deployed context is 
consistent with the defined context in the project YAML file. For example, if you deploy a context instance and later
change the definition of the context in the project YAML file then the running instance will no longer reflect the definition.
In this case you may choose to update the deployed instance using the `agc context deploy` command.

Status will only be shown for contexts for the current user in the current AWS region for the current project. To show
contexts for another project, issue the command from that project's home folder (or subfolder). To display contexts for
another AWS region, you can use a different AWS CLI profile or set the `AWS_PROFILE` environment variable to the 
desired region (e.g `export AWS_REGION=us-west-2`).

{{% alert title="Warning" color="warning" %}}
Because the `status` command will only show contexts that are listed in the project YAML you should take care to `destroy`
any running contexts before deleting them from the project YAML.
{{% /alert %}}

## Costs

Infrastructure deployed for a context is tagged with the context name as well as username and project name. These tags
can be used with AWS CostExplorer to identify the costs associated with running contexts.

A deployed context will incur charges based on the resources being used by the context. If a workflow is running this
will include compute costs for running the workflow tasks but some contexts may include infrastructure that is always
"on" and will incur costs even when no workflow is running. If you no longer need a context we recommend pausing or
destroying it.

If `requestSpotInstances` is true, the context will use spot instances for compute tasks. The context will set the max
price to 100% although if the current price is lower you will pay the lower price. Note that even at 100% spot instances
can still be interrupted if total demand for on demand instances in an availability zone exceeds the available pool. For
full details see [Spot Instance Interruptions](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-interruptions.html) 
and [EC2 Spot Pricing](https://aws.amazon.com/ec2/spot/pricing/).

### Tags

All context infrastructure is [tagged]( {{< relref "namespaces#tags" >}} ) with the context name, username and project name. These tags may be used to help
differentiate costs.

## Technical Details

Context infrastructure is defined as code as [AWS CDK](https://aws.amazon.com/cdk/) apps. For examples, take a look at the `packages/cdk` folder. When 
deployed a context will produce one or more stacks in Cloudformation. Details can be viewed in the Cloudformation console
or with the AWS CLI.

A context includes an endpoint compliant with the [GA4GH WES API](https://ga4gh.github.io/workflow-execution-service-schemas/docs/). This API is how Amazon Genomics CLI submits workflows to the context. The
context also contains one or more workflow engines. These may either be deployed as long-running services as is the case
with Cromwell or as "head" jobs that are responsible for a single workflow, as is the case for NextFlow. Engines run as
"head" jobs are started and stopped on demand thereby saving resources.

### Updating Launch Templates

Changes to EC2 LaunchTemplates in CDK result in a new LaunchTemplate version when the infrastructure is updated. Currently,
CDK is unable to also update the default version of the template. In addition, any existing AWS Batch Compute Environments
will not be updated to use the new LaunchTemplate version. Because of this, whenever a LaunchTemplate is updated in CDK
code we recommend destroying any relevant running contexts and redeploying them. An update will *NOT* be sufficient.
