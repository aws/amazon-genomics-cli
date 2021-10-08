---
title: "Walk through"
date: 2021-09-09T10:16:57-04:00
draft: false
weight: 10
description: >
    Demonstrates installation and using the essential functions of Amazon Genomics CLI
---

## Prerequisites

Ensure you have completed the [prerequisites]( {{< relref "../Getting started/prerequisites" >}} ) before beginning.


## Download and install Amazon Genomics CLI

Download the Amazon Genomics CLI according to the [installation]( {{< relref "../Getting started/installation" >}} )
instructions.

## Setup

Ensure you have initialized your account and created a username by following the [setup]( {{< relref "../Getting started/setup" >}} )
instructions.

## Initialize a project

Amazon Genomics CLI uses local folders and config files to define projects. Projects contain configuration settings for contexts and workflows (more on these below). To create a new project for running WDL based workflows do the following:

```shell
mkdir myproject
cd myproject
agc project init myproject --workflow-type wdl
```

> NOTE: for a Nextflow based project you can substitute `--workflow-type wdl` with `---workflow-type nextflow`.

Projects may have workflows from different languages, so the `--workflow-type` flag is simply to provide the stub for
an initial workflow engine.

This will create a config file called  `agc-project.yaml` with the following contents:

```yaml
name: myproject
schemaVersion: 1
contexts:
    ctx1:
        engines:
            - type: wdl
              engine: cromwell
```

This config file will be used to define aspects of the project - e.g. the contexts and named workflows the project uses. For a
more representative project config, look at the projects in `~/agc/examples`. Unless otherwise stated, command line activities for the remainder of this document will assume they are run from within the `~/agc/examples/demo-wdl-project/` project folder.

## Contexts

Amazon Genomics CLI uses a concept called “contexts” to run workflows. Contexts encapsulate and automate time-consuming tasks 
like configuring and deploying workflow engines, creating data access policies, and tuning compute clusters for operation at scale. 
In the `demo-wdl-project` folder, after running the following:

```shell
agc context list
```

You should see something like:

```
2021-09-22T01:15:41Z 𝒊  Listing contexts.
CONTEXTNAME    myContext
CONTEXTNAME    spotCtx
```

In this project there are two contexts, one configured to run with On-Demand instances (myContext), and one configured to use SPOT instances (spotCtx).

You need to have a context running to be able to run workflows. To deploy the context `myContext` in the demo-wdl-project, run:

```shell
agc context deploy -c myContext
```

This will take 10-15min to complete.

If you have more than one context configured, and want to deploy them all at once, you can run:

```shell
agc context deploy --all
```

Contexts have read-write access to a context specific prefix in the S3 bucket Amazon Genomics CLI creates during account activation. You can check this for the `myContext` context with:

```shell
agc context describe -c myContext
```

You should see something like:

```
CONTEXT    myContext    false    STARTED
OUTPUTLOCATION    s3://agc-123456789012-us-east-2/project/Demo/userid/xxxxxxxxJKP3z/context/myContext
```

You can add more data locations using the `data` section of the `agc-project.yaml` config file. All contexts will have an 
appropriate access policy created for the data locations listed when they are deployed. For example, the following config adds three public buckets from the Registry of Open Data on AWS:

```yaml
name: myproject
data:
  - location: s3://broad-references
    readOnly: true
  - location: s3://gatk-test-data
    readOnly: true
  - location: s3://1000genomes
    readOnly: true
```

Note, you need to redeploy any running contexts to update their access to data locations. Do this by simply (re)running.

```shell
agc context deploy -c myContext
```

Contexts also define what types of compute your workflow will run on - i.e. if you want to run workflows using SPOT or On-demand instances. 
By default, contexts use On-demand instances. The configuration for a context that uses SPOT instances looks like the following:

```yaml
contexts:
  # The spot context uses EC2 spot instances which are usually cheaper but may be interrupted
  spotCtx:
    requestSpotInstances: true
```

You can also explicitly specify what instance types contexts will be able to use for workflow jobs. By default, Amazon Genomics CLI 
will use a select set of instance types optimized for running genomics workflow jobs that balance data I/O performance and 
mitigation of workflow failure due to SPOT reclamation. In short, Amazon Genomics CLI uses AWS Batch for job execution 
and selects instance types based on the requirements of submitted jobs, up to `4xlarge` instance types. If you have a use case
that requires a specific set of instance types, you can define them with something like:

```yaml
contexts:
  specialCtx:
    instanceTypes:
      - c5
      - m5
      - r5
```

The above will create a context called `specialCtx` that will use any size of instances in the C5, M5, and R5 instance families.
Contexts are elastic with a minimum vCPU capacity of 0 and a maximum of 256. When all vCPUs are allocated to jobs, further
tasks will be queued.

Contexts also launch an engine for specific workflow types. You can have one engine per context and, currently, engines for WDL and Nextflow are supported.

A contexts configured with WDL and Nextflow engines respectively look like:

```yaml
contexts:
  wdlContext:
    engines:
      - type: wdl
        engine: cromwell

  nfContext:
    engines:
      - type: nextflow
        engine: nextflow
```

## Workflows

### **Add a workflow**

Bioinformatics workflows are written in languages like WDL and Nextflow in either single script files, or in packages of 
multiple files (e.g. when there are multiple related workflows that leverage reusable elements). Currently, Amazon Genomics CLI supports both WDL and Nextflow. 
To learn more about WDL workflows, we suggest resources like the [OpenWDL - Learn WDL course](https://github.com/openwdl/learn-wdl). To learn more about Nextflow workflows, we suggest [Nextflow’s documentation](https://nextflow.io/docs/latest/index.html) and [NF-Core](https://nf-co.re/).

For clarity, we’ll refer to these workflow script files as “workflow definitions”. A “workflow specification” for Amazon Genomics CLI references workflow definitions and combines it with additional metadata, 
like the workflow language the definition is written in, which Amazon Genomics CLI will use to execute it on appropriate compute resources.

There is a “hello” workflow definition in the `~/agc/examples/demo-wdl-project/workflows/hello` folder that looks like:

```
version 1.0
workflow hello_agc {
    call hello {}
}
task hello {
    command { echo "Hello Amazon Genomics CLI!" }
    runtime {
        docker: "ubuntu:latest"
    }
    output { String out = read_string( stdout() ) }
}
```

The workflow specification for this workflow in the project config looks like:

```yaml
workflows:
  hello:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/hello.wdl
```

Here the workflow is expected to conform to the `WDL-1.0` specification. A specification for a “hello” workflow written in Nextflow DSL1 would look like:

```yaml
workflows:
  hello:
    type:
      language: nextflow
      version: 1.0
    sourceURL: workflows/hello
```

For Nextflow DSL2 workflows set `type.version` to `dsl2`.

NOTE: When referring to local workflow definitions, `sourceURL` must either be a full absolute path or a path relative to the `agc-project.yaml` file. Path expansion is currently not supported.

You can quickly get a list of available configured workflows with:

```shell
agc workflow list
```

For the `demo-wdl-project`, this should return something like:

```
2021-09-02T05:14:47Z 𝒊  Listing workflows.
WORKFLOWNAME    haplotype
WORKFLOWNAME    hello
WORKFLOWNAME    read
WORKFLOWNAME    words-with-vowels
```

The `hello` workflow specification points to a single-file workflow. Workflows can also be directories. For example, if you have a workflow that looks like:

```
workflows/hello-dir
|-- inputs.json
`-- main.wdl
```

The workflow specification for the workflow above would simply point to the parent directory:

```yaml
workflows:
  hello-dir-abs:
    type:
        language: wdl
        version: 1.0
    sourceURL: /abspath/to/hello-dir
  hello-dir-rel:
    type:
        language: wdl
        version: 1.0
    sourceURL: relpath/to/hello-dir
```

In this case, your workflow must be named `main.<workflow-type>` - e.g. `main.wdl`

You can also provide a `MANIFEST.json` file that points to a specific workflow file to run. If you have a folder like:

```
workflows/hello-manifest/
|-- MANIFEST.json
|-- hello.wdl
|-- inputs.json
`-- options.json
```

The `MANIFEST.json` file would be:

```json5
{
    "mainWorkflowURL": "hello.wdl",
    "inputFileURLs": [
        "inputs.json"
    ]
}
```

At minimum, MANIFEST files *must* have a `mainWorkflowURL` property which is a relative path to the workflow file in its parent directory.

Workflows can also be from remote sources like GitHub:

```yaml
workflows:
  remote:
    type:
        language: wdl
        version: 1.0  # this is the WDL spec version
    sourceURL: https://raw.githubusercontent.com/openwdl/learn-wdl/master/1_script_examples/1_hello_worlds/1_hello/hello.wdl
```

> NOTE: remote sourceURLs for Nextflow workflows can be Git repo URLs like: https://github.com/nextflow-io/rnaseq-nf.git

### Running a workflow

To run a workflow you need a running context. See the section on contexts above if you need to start one. To run the “hello” workflow in the “myContext” context, run:

```shell
agc workflow run hello --context myContext
```

If you have another context in your project, for example one named “test”, you can run the “hello” workflow there with:

```shell
agc workflow run hello --context test
```

If your workflow was successfully submitted you should get something like:

```
2021-08-04T23:01:37Z 𝒊  Running workflow. Workflow name: 'hello', Arguments: '', Context: 'myContext'
"06604478-0897-462a-9ad1-47dd5c5717ca"
```

The last line is the workflow run id. You use this id to reference a specific workflow execution.

Running workflows is an asynchronous process. After submitting a workflow from the CLI, it is handled entirely in the cloud. 
You can now close your terminal session if needed. The workflow will still continue to run. You can also run multiple workflows at a time. The underlying compute resources will automatically scale. Try running multiple instances of the “hello” workflow at once.

You can check the status of all running workflows with:

```shell
agc workflow status
```

You should see something like this:

```
WORKFLOWINSTANCE    myContext    66826672-778e-449d-8f28-2274d5b09f05    true    COMPLETE    2021-09-10T21:57:37Z    hello
```

By default, the `workflow status` command will show the state of all workflows across all running contexts.

To show only the status of workflow instances of a specific workflow you can use:

```shell
agc workflow status -n <workflow-name>
```

To show only the status of workflows instances in a specific context you can use:

```shell
agc workflow status -c <context-name>
```

If you want to check the status of a specific workflow you can do so by referencing the workflow execution by its run id:

```shell
agc workflow status -r <workflow-run-id>
```

If you need to stop a running workflow instance, run:

```shell
agc workflow stop <workflow-run-id>
```

### Using workflow inputs

You can provide runtime inputs to workflows at the command line. For example, the `demo-wdl-project` has a workflow named `read` that requires reading a data file.
The specification of read looks like:

```
version 1.0
workflow ReadFile {
    input {
        File input_file
    }
    call read_file { input: input_file = input_file }
}

task read_file {
    input {
        File input_file
    }
    String content = read_string(input_file)

    command {
        echo '~{content}'
    }
    runtime {
        docker: "ubuntu:latest"
        memory: "4G"
    }

    output { String out = read_string( stdout() ) }
}
```

You can create an input file locally for this workflow:

```shell
mkdir inputs
echo "this is some data" > inputs/data.txt
cat << EOF > inputs/read.inputs.json
{"ReadFile.input_file": "data.txt"}
EOF
```

Finally, you would submit the workflow with its corresponding inputs file with:

```shell
agc workflow run read --args inputs/read.inputs.json
```

Amazon Genomics CLI will scan the file provided to `--args` for local paths, sync those files to S3, and rewrite the 
inputs file in transit to point to the appropriate S3 locations. Paths in the `*.inputs.json` file provided as `--args` 
are referenced relative to the `*.inputs.json` file.

### Accessing workflow results

Workflow results are written to an S3 bucket specified or created by Amazon Genomics CLI during account activation. See the section on account activation above for more details. You can list or retrieve the S3 URI for the bucket with:

```shell
AGC_BUCKET=$(aws ssm get-parameter \
    --name /agc/_common/bucket \
    --query 'Parameter.Value' \
    --output text)
```

and then use `aws s3` commands to explore and retrieve data from the bucket. For example, to list the bucket contents:

```shell
aws s3 ls $AGC_BUCKET
```

You should see something like:

```
PRE project/
PRE scripts/
```

Data for multiple projects are kept in `project/<project-name>` prefixes. Looking into one you should see:

```
PRE cromwell-execution/
PRE workflow/
```


The `cromwell-execution` prefix is specific to the engine Amazon Genomics CLI uses to run WDL workflows. 
Workflow results will be in `cromwell-execution` partitioned by workflow name, workflow run id, and task name. The `workflow` prefix is where named workflows are cached when you run workflows definitions stored in your local environment.

### Accessing workflow logs

You can get a summary of the log information for a workflow as follows:

```shell
agc logs workflow <workflow-name>
```

This will return the logs for all runs of the workflow. If you just want the logs for a specific workflow run, you can use:

```shell
agc logs workflow <workflow-name> -r <workflow-instance-id>
```

This will print out the `stdout` generated by each workflow task.

For the hello workflow above, this would look like:

```
Fri, 10 Sep 2021 22:00:04 +0000    download: s3://agc-123456789012-us-east-2/scripts/3e129a27928c192b7804501aabdfc29e to tmp/tmp.BqCb2iaae/batch-file-temp
Fri, 10 Sep 2021 22:00:04 +0000    *** LOCALIZING INPUTS ***
Fri, 10 Sep 2021 22:00:05 +0000    download: s3://agc-123456789012-us-east-2/project/Demo/userid/xxxxxxxxJKP3z/context/myContext/cromwell-execution/hello_agc/66826672-778e-449d-8f28-2274d5b09f05/call-hello/script to agc-123456789012-us-east-2/project/Demo/userid/xxxxxxxxJKP3z/context/myContext/cromwell-execution/hello_agc/66826672-778e-449d-8f28-2274d5b09f05/call-hello/script
Fri, 10 Sep 2021 22:00:05 +0000    *** COMPLETED LOCALIZATION ***
Fri, 10 Sep 2021 22:00:05 +0000    Hello Amazon Genomics CLI!
Fri, 10 Sep 2021 22:00:05 +0000    *** DELOCALIZING OUTPUTS ***
Fri, 10 Sep 2021 22:00:05 +0000    upload: ./hello-rc.txt to s3://agc-123456789012-us-east-2/project/Demo/userid/xxxxxxxxJKP3z/context/myContext/cromwell-execution/hello_agc/66826672-778e-449d-8f28-2274d5b09f05/call-hello/hello-rc.txt
Fri, 10 Sep 2021 22:00:06 +0000    upload: ./hello-stderr.log to s3://agc-123456789012-us-east-2/project/Demo/userid/xxxxxxxxJKP3z/context/myContext/cromwell-execution/hello_agc/66826672-778e-449d-8f28-2274d5b09f05/call-hello/hello-stderr.log
Fri, 10 Sep 2021 22:00:06 +0000    upload: ./hello-stdout.log to s3://agc-123456789012-us-east-2/project/Demo/userid/xxxxxxxxJKP3z/context/myContext/cromwell-execution/hello_agc/66826672-778e-449d-8f28-2274d5b09f05/call-hello/hello-stdout.log
Fri, 10 Sep 2021 22:00:06 +0000    *** COMPLETED DELOCALIZATION ***
```

If your workflow fails, useful debug information is typically reported by the workflow engine logs. These are unique per context. To get those for a context named `myContext`, you would run:

```shell
agc logs engine --context myContext
```

You should get something like:

```
Fri, 10 Sep 2021 23:40:49 +0000    2021-09-10 23:40:49,421 cromwell-system-akka.dispatchers.api-dispatcher-175 INFO  - WDL (1.0) workflow 1473f547-85d8-4402-adfc-e741b7df69f2 submitted
Fri, 10 Sep 2021 23:40:52 +0000    2021-09-10 23:40:52,711 cromwell-system-akka.dispatchers.engine-dispatcher-30 INFO  - 1 new workflows fetched by cromid-2054603: 1473f547-85d8-4402-adfc-e741b7df69f2
Fri, 10 Sep 2021 23:40:52 +0000    2021-09-10 23:40:52,712 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - WorkflowManagerActor: Starting workflow UUID(1473f547-85d8-4402-adfc-e741b7df69f2)
Fri, 10 Sep 2021 23:40:52 +0000    2021-09-10 23:40:52,712 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - WorkflowManagerActor: Successfully started WorkflowActor-1473f547-85d8-4402-adfc-e741b7df69f2
Fri, 10 Sep 2021 23:40:52 +0000    2021-09-10 23:40:52,712 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - Retrieved 1 workflows from the WorkflowStoreActor
Fri, 10 Sep 2021 23:40:52 +0000    2021-09-10 23:40:52,716 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - MaterializeWorkflowDescriptorActor [UUID(1473f547)]: Parsing workflow as WDL 1.0
Fri, 10 Sep 2021 23:40:52 +0000    2021-09-10 23:40:52,721 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - MaterializeWorkflowDescriptorActor [UUID(1473f547)]: Call-to-Backend assignments: hello_agc.hello -> AWSBATCH
Fri, 10 Sep 2021 23:40:52 +0000    2021-09-10 23:40:52,722  WARN  - Unrecognized configuration key(s) for AwsBatch: auth, numCreateDefinitionAttempts, filesystems.s3.duplication-strategy, numSubmitAttempts, default-runtime-attributes.scriptBucketName
Fri, 10 Sep 2021 23:40:53 +0000    2021-09-10 23:40:53,741 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - WorkflowExecutionActor-1473f547-85d8-4402-adfc-e741b7df69f2 [UUID(1473f547)]: Starting hello_agc.hello
Fri, 10 Sep 2021 23:40:54 +0000    2021-09-10 23:40:54,030 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - Assigned new job execution tokens to the following groups: 1473f547: 1
Fri, 10 Sep 2021 23:40:55 +0000    2021-09-10 23:40:55,501 cromwell-system-akka.dispatchers.engine-dispatcher-4 INFO  - 1473f547-85d8-4402-adfc-e741b7df69f2-EngineJobExecutionActor-hello_agc.hello:NA:1 [UUID(1473f547)]: Call cache hit process had 0 total hit failures before completing successfully
Fri, 10 Sep 2021 23:40:56 +0000    2021-09-10 23:40:56,842 cromwell-system-akka.dispatchers.engine-dispatcher-31 INFO  - WorkflowExecutionActor-1473f547-85d8-4402-adfc-e741b7df69f2 [UUID(1473f547)]: Job results retrieved (CallCached): 'hello_agc.hello' (scatter index: None, attempt 1)
Fri, 10 Sep 2021 23:40:57 +0000    2021-09-10 23:40:57,820 cromwell-system-akka.dispatchers.engine-dispatcher-4 INFO  - WorkflowExecutionActor-1473f547-85d8-4402-adfc-e741b7df69f2 [UUID(1473f547)]: Workflow hello_agc complete. Final Outputs:
Fri, 10 Sep 2021 23:40:57 +0000    {
Fri, 10 Sep 2021 23:40:57 +0000      "hello_agc.hello.out": "Hello Amazon Genomics CLI!"
Fri, 10 Sep 2021 23:40:57 +0000    }
Fri, 10 Sep 2021 23:40:59 +0000    2021-09-10 23:40:59,826 cromwell-system-akka.dispatchers.engine-dispatcher-14 INFO  - WorkflowManagerActor: Workflow actor for 1473f547-85d8-4402-adfc-e741b7df69f2 completed with status 'Succeeded'. The workflow will be removed from the workflow store.
```

You can filter logs with the `--filter` flag. The filter syntax adheres to [CloudWatch's filter and pattern syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).
For example, the following will give you all error logs from the workflow engine:

```shell
agc logs engine --context myContext --filter ERROR
```

### Additional workflow examples

The Amazon Genomics CLI installation also includes a set of typical genomics workflows for raw data processing, germline variant discovery, and joint genotyping based on [GATK Best Practices](https://gatk.broadinstitute.org/hc/en-us), developed by the [Broad Institute](https://www.broadinstitute.org/). More information on how these workflows work is available in the [GATK Workflows Github repository](https://github.com/gatk-workflows).

You can find these in:

```shell
~/agc/examples/gatk-best-practices-project
```

These workflows come pre-packaged with `MANIFEST.json` files that specify example input data available publicly in the [AWS Registry of Open Data](https://registry.opendata.aws/).

Note: these workflows take between 5 min to ~3hrs to complete.

## Cleanup

When you are done running workflows, it is recommended you stop all cloud resources to save costs.

Stop a context with:

```shell
agc context destroy -c <context-name>
```

This will destroy all compute resources in a context, but retain any data in S3. If you want to destroy all your running contexts at once, you can use:

```shell
agc context destroy --all
```

Note, you will not be able to destroy a context that has a running workflow. Workflows will need to complete on their own or stopped before you can destroy the context.

If you want stop using Amazon Genomics CLI in your AWS account entirely, you need to deactivate it:

```shell
agc account deactivate
```

This will remove Amazon Genomics CLI’s core infrastructure. If Amazon Genomics CLI created a VPC as part of the `activate` process, it will be **removed**. If Amazon Genomics CLI created an S3 bucket for you, it will be **retained**.

To uninstall Amazon Genomics CLI from your local machine, run the following command:

```shell
./agc/uninstall.sh
```

> Note uninstalling the CLI will not remove any resources or persistent data from your AWS account.
