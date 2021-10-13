---
title: "Hello world"
date: 2021-09-07T14:01:59-04:00
draft: false
weight: 40
---

When you install Amazon Genomics CLI it will create a folder named `agc`. Inside there is an `examples/demo-project` folder containing an `agc-project.yaml`
with some demo projects including a simple "hello world" workflow.

## Running Hello World

1. Ensure you have met the [prerequisites]( {{< relref "prerequisites" >}} ) and [installed]( {{< relref "installation" >}} ) Amazon Genomics CLI
2. Ensure you have followed the [activation]( {{< relref "setup" >}} ) steps
3. `cd ~/agc/examples/demo-wdl-project`
4. `agc context deploy --context myContext`, this step takes approximately 5 minutes to deploy the infrastructure
5. `agc workflow run hello --context myContext`, take note of the returned workflow instance ID.
6. Check on the status of the workflow `agc workflow status -r <workflow-instance-id>`. Initially you will see status like `SUBMITTED` but after the elastic compute resources have been spun up and the workflow runs you should see something like the following: `WORKFLOWINSTANCE    ctx1    9ff7600a-6d6e-4bda-9ab6-c615f5d90734    COMPLETE    2021-09-01T20:17:49Z`

Congratulations! You have just run your first workflow in the cloud using Amazon Genomics CLI! At this point you can run additional workflows, including submitting several instances of the "hello world" workflow.
The elastic compute resources will expand and contract as necessary to accommodate the backlog of submitted workflows.

## Reviewing the Results

Workflow results are written to an S3 bucket specified or created by Amazon Genomics CLI during account activation. 
You can list or retrieve the S3 URI for the bucket with:

```shell
AGC_BUCKET=$(aws ssm get-parameter \
    --name /agc/_common/bucket \
    --query 'Parameter.Value' \
    --output text)
```

and then use `aws s3` commands to explore and retrieve data from the bucket. Workflow output will be in the 
`s3://agc-<account-num>-<region>/project/<project-name>/userid/<user-id>/context/<context-name>/workflow/<workflow-name>/`
path. The rest of the path depends on the engine used to run the workflow. For Cromwell it will continue with:
`.../cromwell-execution/<wdl-wf-name>/<workflow-run-id>/<task-name>`

If a workflow declares outputs then you may obtain these using the command:

```shell
agc workflow output <workflow_run_id>
```

You should see a response similar to:

```shell
OUTPUT	id	6cc6f742-dc87-4649-b319-1af45c4c09c6
OUTPUT	outputs.hello_agc.hello.out	Hello Amazon Genomics CLI!
```

You can also obtain task logs for a workflow using the following form `agc logs workflow <workflow-name> -r <instance-id>`.
>Note, if the workflow did not actually run any tasks due to call caching then there will be no output from this command.

## Cleaning Up

Once you are done with `myContext` you can destroy it with:

```shell
agc context destroy -c myContext
```

This will remove the cloud resources associated with the named context, but will keep any S3 outputs and CloudWatch logs.

If you want stop using Amazon Genomics CLI in your AWS account entirely, you need to deactivate it:

```shell
agc account deactivate
```

This will remove Amazon Genomics CLI’s core infrastructure. If Amazon Genomics CLI created a VPC as part of the activate process, it will be *removed*. If Amazon Genomics CLI created an S3 bucket for you, it will be *retained*.

To uninstall Amazon Genomics CLI from your local machine, run the following command:

```shell
./agc/uninstall.sh

```

Note uninstalling the CLI will *not* remove any resources or persistent data from your AWS account.

## Next Steps

* Familiarize yourself with [Amazon Genomics CLI Concepts]( {{< ref "Concepts" >}} )
* Try some [tutorials]( {{< ref "Tutorials" >}} )
