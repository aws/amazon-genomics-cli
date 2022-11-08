---
title: "Contexts"
date: 2021-08-31T17:26:49-04:00
draft: false
weight: 20
description: >
    Contexts are the set of cloud resources used to run a workflow
---
## What is a Context?

A context is a set of cloud resources. Amazon Genomics CLI runs [workflows]( {{< relref "../workflows" >}} ) in a context. A deployed context will include an 
[engine]( {{< relref "../engines" >}}) that can interpret
and manage the running of a workflow along with compute resources that will run the individual tasks of the workflow. The
deployed context will also contain any resources needed by the engine or compute resources including any security, permissions
and [logging]( {{< relref "../logs" >}} ) capabilities. Deployed contexts are [namespaced]( {{< relref "../namespaces" >}}) based on the user, project and context name so that resources
are isolated, preventing collisions.

When a workflow is run the user will decide which context will run it. For example, you might choose to submit a workflow
to a context that uses "Spot priced" resources or one that uses "On Demand" priced resources.

When deployed context resources that require a VPC will be deployed into the VPC that was specified when the [account]( {{< relref "../accounts" >}} ) was
activated.

## How is a Context Defined?

A context is defined in the YAML file that defines the [project]( {{< relref "../projects" >}} ). A project has at least one context but may have many.
Contexts must have unique names and are defined as YAML maps.

A context may request use of [Spot priced](https://aws.amazon.com/ec2/spot/pricing/) compute resources with `requestSpotInstances: true`. The default value is `false`.

A context must define an array of one or more `engines`. Each engine definition must specify the workflow language that it 
will interpret. For each language Amazon Genomics CLI has a default engine however, users may specify the exact engine in the `engine`
parameter.

## General Architecture of a Context

The exact architecture of a context will depend on the context properties described below and defined in their `agc-project.yaml`. However, the architecture deployed on execution of `agc context deploy` is shown in the following diagram:

![Image of the general architecture of a context](ContextGeneralArchitecture.png "General Architecture of a Context")

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

### Public Subnets

In the interest of saving money, in particular if you intend to have the AGC stack deployed for a long period, you may choose to deploy in "public subnet" mode.
To do this, you must first set up the core stack using `aws configure --usePublicSubnets`, which will disable the creation of the NAT gateway and VPC endpoints which present an ongoing cost unrelated to your use of compute resources.
After you have done this, you must also set `usePublicSubnets: true` in all contexts you use:
```yaml
contexts:
  someCtx:
    usePublicSubnets: true
    engines:
      - type: nextflow
        engine: nextflow
```

This ensures that the AWS batch instances are deployed into a public subnet, which has no additional cost associated with it.
However note that while these instances are given a security group that will block all incoming traffic, this is not as secure as using the default private subnet mode.

## Context Commands

A full reference of context commands is [here]( {{< relref "../../Reference/agc_context" >}} )

### `describe`

The command `agc context describe <context-name> [flags]` will describe the named context as defined in the project YAML
as well as other relevant account information.

### `list`

The command `agc context list [flags]` will list the names of all contexts defined in the project YAML file along with the name of the engine used by the context.

### `deploy`

The command `agc context deploy <context-name> [flags]` is used to deploy the cloud infrastructure required by the context.
If the context is already running the existing infrastructure will be updated to reflect changes in project YAML. For example
if you added another `data` definition in your project and run `agc context deploy <context-name>` then the deployed context
will be updated to allow access to the new data.

All contexts defined in the project YAML can be deployed or updated using the `--all` flag.

Individually named contexts can be deployed or updated as positional arguments. For example: `agc context deploy ctx1 ctx2`
will deploy the contexts `ctx1` and `ctx2`.

The inclusion of the `--verbose` flag will show the full CloudFormation output of the context deployment.

### `destroy`

A contexts cloud resources can be "destroyed" using the `agc context destroy <context-name>` command. This will remove any 
infrastructure artifacts associated with the context unless they are defined as being retained. Typically, things like logs
and workflow outputs on S3 are retained when a context is destroyed.

All deployed contexts can be destroyed using the `--all` flag.

Multiple contexts can be destroyed in a single command using positional arguments. For example: `agc context destroy ctx1 ctx2`
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

### Ongoing Costs

Until a context is destroyed resources that are deployed can incur ongoing costs even if a workflow is not running. The
exact costs depend on the configuration of the context.

Amazon Genomics CLI version 1.0.1 and earlier used an AWS Fargate based WES service for each deployed context. The service
uses 0.5 vCPU, 4 GB memory and 20 GB base instance storage. Fargate pricing varies by region and is detailed [here](https://aws.amazon.com/fargate/pricing/).
The estimated cost is available via [this link](https://calculator.aws/#/estimate?id=9a67ba7845199cf108d85ae0f9b8176253266005)

After version 1.0.1, the WES endpoints deployed by Amazon Genomics CLI are implemented with AWS Lambda and therefore use
a [pricing model](https://aws.amazon.com/lambda/pricing/) based on invocations.

Contexts using a Cromwell engine run an additional AWS Fargate service for the engine with 2 vCPU, 16 GB RAM and 20 GB of
base storage. Additionally, Cromwell is deployed with a standard EFS volume for storage of metadata. EFS [costs](https://aws.amazon.com/efs/pricing/) are volume based. While
relatively small the amount of metadata will expand as more workflows are run. The volume is destroyed when the context is destroyed. An estimated
cost for both components is available via [this link](https://calculator.aws/#/estimate?id=8ccc606c1b267e2933a6d683c0b98fcf11e4cbab)

Contexts using the "miniwdl" or "snakemake" engines use EFS volumes as scratch space for workflow intermediates, caches and temporary files. Because many genomics
workflows can accumulate several GB of intermediates per run we recommend destroying these contexts when not in use. An estimated cost assuming a
total of 500 GB of workflow artifacts is available via [this link](https://calculator.aws/#/estimate?id=4d19b43aa86fcc3af199c425bfcc55193592cbb4)

Refer to the [public subnets section](#public-subnets) if you are concerned about reducing these ongoing costs.

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
