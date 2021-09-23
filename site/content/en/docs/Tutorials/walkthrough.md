---
title: "Walk through"
date: 2021-09-09T10:16:57-04:00
draft: false
weight: 10
description: >
    Demonstrates installation and using the essential functions of AGC
---

## Prerequisites

Ensure you have completed the [prerequisites]( {{< relref "../Getting started/prerequisites" >}} ) before beginning.


## Download and install Amazon Genomics CLI

Download the Amazon Genomics CLI according to the [installation]( {{< relref "../Getting started/installation" >}} )
instructions.

## Setup

Ensure you have initialized your account and created a username by following the [setup]( {{< relref "../Getting started/setup" >}} )
instructions.

## Create a Project

Amazon Genomics CLI uses local folders and config files to define projects. Projects contain configuration settings for contexts and workflows (more on these below). To create a new project do the following:

```shell
mkdir myproject
cd myproject
agc project init myproject
```

This will create a config file called  `agc-project.yaml` with the following content:

```yaml
name: myproject
contexts:
ctx1: {}
```

This config file will be used to define aspects of the project - e.g. contexts and named workflows the project uses. 
For more representative project YAML files, look at the example projects in `~/agc/examples` which are created when AGC 
is installed.

## Contexts


Amazon Genomics CLI uses a concept called “contexts” to run workflows. Contexts encapsulate and automate time-consuming tasks like
configuring and deploying workflow engines, creating data access policies, and tuning compute clusters for operation at
scale. When a project is first initialized, the resultant project config will include a configuration for a context called 
`ctx1`. In the `myproject` folder you created above, after running the following:

```shell
agc context list
```

You should see something like:

```
2021-07-31T03:44:51Z 𝒊  Listing contexts.
CONTEXTNAME    ctx1
```

You can customize this context as needed. In the ~/agc/examples/demo-project this has been reconfigured and renamed myContext:

```shell
cd ~/agc/examples/demo-project
agc context list
```

```
2021-09-02T04:12:33Z 𝒊  Listing contexts.
CONTEXTNAME    myContext
```

You need to have a context running to be able to run workflows. To deploy the context `myContext` in the demo-project, run:

```shell
agc context deploy myContext
```

This will take 5-10min to complete.

If you have more than one context configured, and want to deploy them all at once, you can run:

```
agc context deploy --all
```

Contexts have read-write access to a context specific prefix in the S3 bucket Amazon Genomics CLI creates during account activation. You can check this for the `myContext` context with:

```shell
agc context describe myContext
```

You should see something like:

```
CONTEXT    myContext    false    NOT_STARTED
OUTPUTLOCATION    s3://agc-733263974272-us-east-2/project/Demo/userid/pwymingJKP3z/context/myContext
```

You can add more data locations using the data section of the `agc-project.yaml` config file. All contexts will have an appropriate access policy created for the data locations listed when they are deployed. For example, the following config adds three public buckets from the Registry of Open Data on AWS:

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
agc context deploy myContext
```

Contexts also define what types of compute your workflow will run on - i.e. if you want to run workflows using SPOT or On-demand instances. By default, contexts use On-demand instances. The configuration for a context that uses SPOT instances looks like the following:

```yaml
contexts:
# The spot context uses EC2 spot instances which are usually cheaper but may be interrupted
spotCtx:
requestSpotInstances: true
```

## Workflows

### **Add a workflow**

Bioinformatics workflows are written in languages like WDL in either single script files, or in packages of multiple files (e.g. when there are multiple related workflows that leverage reusable elements). Currently WDL is the only language that Amazon Genomics CLI supports. To learn more about writing WDL workflows, we suggest resources like the [OpenWDL - Learn WDL course](https://github.com/openwdl/learn-wdl).

For clarity, we’ll refer to these workflow script files as “workflow definitions”. A “workflow specification” for Amazon Genomics CLI references workflow definitions and combines it with additional metadata, like the workflow language the definition is written in, which Amazon Genomics CLI will use to execute it on appropriate compute resources.

There is a “hello” workflow in the `~/agc/examples/demo-project` folder that looks like:

```
version 1.0
workflow w {
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
    type: wdl
    version: 1.0  # this is the WDL spec version
    sourceURL: /path/to/hello.wdl
```

NOTE: When referring to local workflow definitions, `sourceURL` must either be a full absolute path or a path relative to the `agc-project.yaml` file. Path expansion is currently not supported.

You can quickly get a list of available configured workflows with:

```shell
agc workflow list
```

For the `demo-project`, this should return something like:

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
    type: wdl
    sourceURL: /abspath/to/hello-dir
  hello-dir-rel:
    type: wdl
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

At minimum, MANIFEST files must have a `mainWorkflowURL` property which is a relative path to the workflow file in its parent directory.

Workflows can also be from remote sources like GitHub:

```
workflows:
  remote:
    type: wdl
    sourceURL: https://raw.githubusercontent.com/openwdl/learn-wdl/master/1_script_examples/1_hello_worlds/1_hello/hello.wdl
```

### Running a workflow

To run a workflow you need a running context. See the section on contexts above if you need to start one. To run the “hello” workflow in the “myContext” context, run:

```shell
agc workflow run hello --context myContext
```

If you have another context in your project, for example one named “testCtx”, you can run the “hello” workflow there with:

```shell
agc workflow run hello --context testCtx
```

If your workflow was successfully submitted you should get something like:

```
2021-08-04T23:01:37Z 𝒊  Running workflow. Workflow name: 'hello', Arguments: '', Context: 'myContext'
"06604478-0897-462a-9ad1-47dd5c5717ca"
```

The last line is the workflow execution id. You use this id to reference a specific workflow execution.

Running workflows is an asynchronous process. After submitting a workflow from the CLI, it is handled entirely in the cloud. You can now close your terminal session if needed. The workflow will still continue to run. You can also run multiple workflows at a time. The underlying compute resources will automatically scale. Try running multiple instances of the “hello” workflow at once.

You can check the status of all running workflows with:

```shell
agc workflow status
```

You should see something like this:

```
WORKFLOWINSTANCE    ctx1    9ff7600a-6d6e-4bda-9ab6-c615f5d90734    COMPLETE    2021-09-01T20:17:49Z
```

For more information, you can use:

```shell
agc workflow status -l
```

This will provide extra details like if a workflow completed with an error and workflow execution duration:

```
### TODO UPDATE ###
```

The columns are `duration, execution-id, info, workflow-name, start-time, status`

If you have multiple workflows running simultaneously the above commands will show the state of all of them, as well as workflows that have previously completed.

If you want to check the status of a specific workflow you can do so by referencing the workflow execution by it’s unique id:

```shell
agc workflow status <workflow-instance-id>
```

If you need to stop a running workflow, run:

```shell
agc workflow stop <workflow-instance-id>
```

### Using workflow inputs

You can provide runtime inputs to workflows at the command line. For example, the `demo-project` has a workflow named  `read` that requires reading a data file that looks like:

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

Amazon Genomics CLI will scan the file provided to `--args` for local paths, sync those files to S3, and rewrite the inputs file in transit to point to the appropriate S3 locations. Paths in the `*.inputs.json` file provided as `--args` are referenced relative to the `*.inputs.json` file.

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


The `cromwell-execution` prefix is specific to the engine Amazon Genomics CLI uses to run WDL workflows. Workflow results will be in `cromwell-execution` partitioned by workflow name, workflow execution id, and task name. The `workflow` prefix is where named workflows are cached when you run workflows definitions stored in your local environment.

### Accessing workflow logs

You can get a summary of the log information for a workflow as follows:

```shell
agc logs workflow <workflow-name>
```

This will return the logs for all runs of the workflow. If you just want the logs for a specific workflow run, you can use:

```shell
agc logs workflow <workflow-name> -r <workflow-instance-id>
```

You should see something like this:

```
LogStreamName 1e0fdb8f-2522-46a2-8064-6109a936eebd:
  TaskLogs:
    Name: w.hello
    LogStreamName: 82004336-17eb-4ee1-9d3d-ba4cf3b06b83
    StartTime: 2021-06-16T18:25:08.170Z,    EndTime: 2021-06-16T18:25:53.594Z
    Stdout: s3://agc-<account-id>-<region>/project/<project-name>/cromwell-execution/w/1e0fdb8f-2522-46a2-8064-6109a936eebd/call-hello/hello-stdout.log
    Stderr: s3://agc-<account-id>-<region>/project/<project-name>/cromwell-execution/w/1e0fdb8f-2522-46a2-8064-6109a936eebd/call-hello/hello-stderr.log
    CW log:
      log-group: /aws/batch/job
      log-stream: cromwell_ubuntu_latest72d98123b11b4086634680ff136112993c23763c/default/ebf7361dbf354437bfb9b15053a145f6
    ExitCode: 0
```

For each task – which corresponds to an AWS Batch job that was run – you will be told what stdout/stderr files exist in the
S3 output bucket and what Cloudwatch log streams exist.

To retrieve the content of a specific Cloudwatch log stream, use the following command:

```shell
agc logs log [<log-stream-group>] <log-stream-name>
```

The values for `log-stream-group` and `log-stream-name` to use are in the `CW log` section of each task in the workflow log. In the example above, these are:

```
log-group: /aws/batch/job
log-stream: cromwell_ubuntu_latest72d98123b11b4086634680ff136112993c23763c/default/ebf7361dbf354437bfb9b15053a145f6
```

Logs for individual tasks look like:

```
Event messages for Cloudwatch log stream {/aws/batch/job, cromwell_ubuntu_latest72d98123b11b4086634680ff136112993c23763c/default/ebf7361dbf354437bfb9b15053a145f6}:
  *** LOCALIZING INPUTS ***
  download: s3://agc-<account-id>-<region>/project/<project-name>/cromwell-execution/w/1e0fdb8f-2522-46a2-8064-6109a936eebd/call-hello/script to phosphate-output-bucket-733263974272-us-west-2-demo/cromwell-execution/w/1e0fdb8f-2522-46a2-8064-6109a936eebd/call-hello/script
  *** COMPLETED LOCALIZATION ***
  Hello Amazon Genomics CLI!
  *** DELOCALIZING OUTPUTS ***
  upload: ./hello-rc.txt to s3://agc-<account-id>-<region>/project/<project-name>/cromwell-execution/w/1e0fdb8f-2522-46a2-8064-6109a936eebd/call-hello/hello-rc.txt
  upload: ./hello-stderr.log to s3://agc-<account-id>-<region>/project/<project-name>/cromwell-execution/w/1e0fdb8f-2522-46a2-8064-6109a936eebd/call-hello/hello-stderr.log
  upload: ./hello-stdout.log to s3://agc-<account-id>-<region>/project/<project-name>/cromwell-execution/w/1e0fdb8f-2522-46a2-8064-6109a936eebd/call-hello/hello-stdout.log
  *** COMPLETED DELOCALIZATION ***
```

If you omit the log-stream-group then the AWS Batch log stream group (/aws/batch/job/) will be assumed.

If your workflow fails, useful debug information is typically reported by the workflow engine logs. These are unique per context. To get those for a context named `myContext`, you would run:

```shell
agc logs engine --context myContext
```

You should get something like:

```
Event messages for Cloudwatch log stream {/agc/myproj/prod/cromwell---11e32d1d-6812-4de0-9581-2c795a7f0831, cromwell/web/d8901182b3154adca986953da11feb52}:
  2021-07-07 04:39:44,705  INFO  - Running with database db.url = jdbc:hsqldb:mem:a8c7c9a2-df8e-442f-ab62-db3ad3267dd5;shutdown=false;hsqldb.tx=mvcc
  2021-07-07 04:39:53,424  INFO  - Running migration RenameWorkflowOptionsInMetadata with a read batch size of 100000 and a write batch size of 100000
  2021-07-07 04:39:53,435  INFO  - [RenameWorkflowOptionsInMetadata] 100%
  2021-07-07 04:39:53,588  INFO  - Running with database db.url = jdbc:hsqldb:mem:3e482536-11e0-439d-8b34-cfc98d8f7810;shutdown=false;hsqldb.tx=mvcc
  2021-07-07 04:39:54,050  WARN  - Unrecognized configuration key(s) for AwsBatch: auth, numCreateDefinitionAttempts, filesystems.s3.duplication-strategy, numSubmitAttempts, default-runtime-attributes.scriptBucketName
  2021-07-07 04:39:54,297  INFO  - Slf4jLogger started
  2021-07-07 04:39:54,528 cromwell-system-akka.dispatchers.engine-dispatcher-4 INFO  - Workflow heartbeat configuration:
  {
    "cromwellId" : "cromid-f52e2fc",
    "heartbeatInterval" : "2 minutes",
    "ttl" : "10 minutes",
    "failureShutdownDuration" : "5 minutes",
    "writeBatchSize" : 10000,
    "writeThreshold" : 10000
  }
  2021-07-07 04:39:54,727 cromwell-system-akka.dispatchers.service-dispatcher-8 INFO  - Metadata summary refreshing every 1 second.
  
[... truncated ...]

  2021-07-07 04:55:41,717 cromwell-system-akka.dispatchers.engine-dispatcher-20 INFO  - WorkflowExecutionActor-f7ec14a3-ef3f-4ab8-8a94-dd3a8f8c0ae5 [UUID(f7ec14a3)]: Workflow w complete. Final Outputs:
  {
    "w.hello.out": "Hello Amazon Genomics CLI!"
  }
  2021-07-07 04:55:41,720 cromwell-system-akka.dispatchers.engine-dispatcher-20 INFO  - WorkflowManagerActor WorkflowActor-f7ec14a3-ef3f-4ab8-8a94-dd3a8f8c0ae5 is in a terminal state: WorkflowSucceededState
  2021-07-07 04:55:41,720 cromwell-system-akka.dispatchers.engine-dispatcher-20 INFO  - WorkflowManagerActor WorkflowActor-cd7b3d2d-6e57-4957-897b-4b092331da96 is in a terminal state: WorkflowSucceededState
```

## Cleanup

When you are done running workflows, it is recommended you stop all cloud resources to save costs.

Stop a context with:

```shell
agc context destroy <context-name>
```

This will destroy all compute resources in a context, but retain any data in S3. If you want to destroy all your running contexts at onece, you can use:

```shell
agc context destroy --all
```

Note, you will not be able to destroy a context that has a running workflow. Workflow will need to complete on their own or stopped before you can destroy the context.

If you want stop using Amazon Genomics CLI in your AWS account entirely, you need to deactivate it:

```shell
agc account deactivate
```

This will remove Amazon Genomics CLI’s core infrastructure. If Amazon Genomics CLI created a VPC as part of the `activate` process, it will be **removed**. If Amazon Genomics CLI created an S3 bucket for you it will be **retained**.

To uninstall Amazon Genomics CLI from your local machine, run the following command:

```shell
./agc/uninstall.sh
```

> Note uninstalling the CLI will not remove any resources or persistent data from your AWS account.
