---
title: "Workflows"
date: 2021-08-31T17:27:16-04:00
draft: false
weight: 40
description: >
    A Workflow is a series of steps or tasks to be executed as part of an analysis.
---
A Workflow is a series of steps or tasks to be executed as part of an analysis. To run a workflow using Amazon Genomics CLI, first you must
have deployed a context with suitable compute resources and with a workflow engine that can interpret the language of
the workflow.

## Specification in Project YAML

In an Amazon Genomics CLI project you can specify multiple workflows in a YAML map. The following example defines four WDL version 1.0
workflows. The `sourceURL` property defines the location of the workflow file. If the location is relative then the
relevant file is assumed to be relative to the location of the project YAML file. Absolute file locations are also possible
although this may reduce the portability of a project if it is intended to be shared. Web URLS are supported as locations
of the workflow definition file.

At this time Amazon Genomics CLI does *not* resolve path aliases so, for example, a `sourceURL` like `~/workflows/worklfow.wdl` is not
supported.

The `type` object declares the `language` of the workflow (eg, wdl, nextflow etc). The run a workflow there must be a
deployed [context]( {{< relref "contexts" >}} ) with a matching language. The `version` property refers to the workflow language version.

```yaml
workflows:
  hello:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/hello.wdl
  read:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/read.wdl
  haplotype:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/haplotypecaller-gvcf-gatk4.wdl
  words-with-vowels:
    type:
      language: wdl
      version: 1.0
    sourceURL: workflows/words-with-vowels.wdl
```

### Multi-file Workflows

Some workflow languages allow for the import of other workflows. To accommodate this, Amazon Genomics CLI supports using a directory as
a source URL. When a directory is supplied as the `sourceURL`, Amazon Genomics CLI uses the following rules to determine the name of
the main workflow file and any supporting files:

1. If the source URL resolves to a single non-zipped file, then the file is assumed to be a workflow file. Dependent resources (if any) are hardcoded in the file and must be resolvable by the Wes adapter or implicitly the workflow engine (e.g the Wes adapter figures out if the engine can resolve them and if not it resolves them itself).
2. The source URL resolves to a zipped file (`.zip`).  The zip may contain a manifest.
   1. If the zip file does *not* contain a file named `MANIFEST.json`:
      1. The zip file must contain one workflow file with the prefix main followed by the conventional suffix for the workflow, e.g. `main.wdl`
      2. Any sub-workflows or tasks referenced by the main workflow must either be in the zip at the appropriate relative path or they must be referenced by URLs that are resolvable by the workflow engine. The WesAdapter may attempt to resolve them for the engine but this is a convenience and not required.
      3. Any variables not defined in the workflows must be provided in an inputs file named with the prefix inputs and the conventional suffix for the workflow engine (e.g `inputs.json`). For workflow engines that support multiple input files an index suffix must be provided (e.g. inputs_a.json or inputs_1.json) if there is more than one inputs file.
      4. A workflow options file may be included and must be named with the options prefix followed by the conventional suffix of the workflow. The WesAdapter may chose to make use of this depending on the context of the workflow engine. It may also choose to pass this to the workflow engine or pass a modified copy to the workflow engine.
   2. If the zip file *does* contain a manifest:
      1. The manifest must contain a parameter called mainWorkflowURL. If it does then the value of the parameter must either be a URL, including the relevant protocol, or the name of a file present in the zip archive. Any subworkflows or tasks imported by the main workflow must either be referenced as URLs in the workflow or be present in the archive as described above.
      2. The manifest may contain an array of URLs to inputs files called inputFileURLs. The WesAdapter must decide if it should resolve these or let the workflow engine resolve them.
      3. The manifest may contain a URL reference to an options files name optionFileURL. The WesAdapter may choose to make use of this depending on the context of the workflow engine. It may also choose to pass this to the workflow engine or pass a modified copy to the workflow engine.
3. If the source URL points to a directory then Amazon Genomics CLI will zip the directory before uploading it. The directory must follow the same conventions stated above for zip files.

The following snippet demonstrates a possible declaration of a multi-file workflow:

```yaml
workflows:
  gatk4-data-processing:
    type:
      langauge: wdl
      version: 1.0
    sourceURL: ./gatk4-data-processing
```

The following snippet demonstrates a valid `MANIFEST.json` file:

```json5
{
  "mainWorkflowURL": "processing-for-variant-discovery-gatk4.wdl",
  "inputFileURLs": [
    "processing-for-variant-discovery-gatk4.hg38.wgs.inputs.json"
  ],
  "optionFileURL": "options.json"
}
```

#### `MANIFEST.json` Structure

The following keys are allowed in the `MANIFEST.json`


| Key               | Required | Purpose                                                                                                                                                                                                                                                                                                                                                                             |
|-------------------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `mainWorkflowURL` | Yes      | Points to the workflow definition that is the main entrypoint of the workflow.                                                                                                                                                                                                                                                                                                      |
| `inputFileURLs`   | No       | An array of URLs to one or more JSON files that define the inputs to the workflow in the format expected by the relevant engine. inputFile URLs are resolved relative to the location of the `MANIFEST.json`.                                                                                                                                                                       |
| `optionFileURL`   | No       | A URL pointing to a JSON file containing engine options applied to a workflow instance. This is only used when engines run in [server mode]( {{< relref "../Concepts/engines#run-mode" >}} ). Options are interpreted by the engine and so must be in the form expected by the engine. The URL is resolved relative to the location of the `MANIFEST.json`.                         |
| `engineOptions`   | No       | A string appended to the command line of the engine's run command. The string may contain any flags or parameters relevant to the engine of the context used to run the workflow. It should not be used to declare inputs (use `inputFileURLS` instead). This parameter is only relevant for engines that run as [head processes]( {{< relref "../Concepts/engines#run-mode" >}} ). |

## Engine Selection

When a workflow is submitted to run, Amazon Genomics CLI will match the workflow type with the map of engines in the context. For example,
if the workflow `type` is `wdl` Amazon Genomics CLI will attempt to identify and engine designated as the engine for that type. There
may only be one engine per type. If no suitable engine is found in the context an error will be reported.

## Workflow Instances

Any defined project workflow can be run multiple times. Each run is called an instance and assigned a unique instance ID.
When referring to a specific run of a workflow you should use the instance ID rather than the workflow name. It is possible
to submit multiple instances of the same workflow and to have these run concurrently.

## Context

All workflows are coordinated by the engine, they are submitted to and executed in the context that is specified at submission time.
The workflow engine decides how the workflow is to be run and the context provides compute resources to run the workflow.

## Commands

A full reference of workflow commands is available [here]( {{< relref "../Reference/agc_workflow" >}} )

### `run`

Invoking `agc workflow run <workflow-name> -c <context-name>` will run the named workflow in a specific context. The
unique ID of that workflow instance run will be returned if the submission is successful.

#### `workflow arguments`

Workflow arguments such as options files can be specified at submission time using the `a` or `--args` flag. For
example:

```shell
agc workflow run my-workflow --args inputs.json
```

If the inputs file references local files, these will be synced with S3 and those files in S3 will be used when the workflow
instance is run.

### `list`

The `agc workflow list` command can be used to list all workflows that are specified in the current project.

### `describe`

The `agc workflow describe <workflow-name>` command will return detailed information about the named workflow based on
the specification in the current project YAML file.

### `status`

To find out the status of workflow instances that are running, or have been run you can use the `agc workflow status` command.
This will display details on 20 recent workflows from the project, to display more, or fewer you can use the `--limit number` flag
where the `number` may be as many as 1000.

To list the status of workflows run or running in a specific context use the `--context-name` flag and provide the name
of one of the contexts of the project.

You may get the status of workflow instances by workflow name using the `--workflow-name` flag.

To display the status of a specific workflow instance you can provide the id of the desired workflow instance with the `--instance-id` flag.

### `stop`

A running workflow *instance* can be stopped at any time using the `agc workflow stop <instance-id>` command. When issued,
Amazon Genomics CLI will look up the appropriate context and engine using the `instance-id` of the workflow and instruct the engine to
stop the workflow. What happens next depends on the actual workflow engine. For example, in the case of the Cromwell WDL
engine, any currently executing tasks will halt, any pending tasks will be removed from the work queue and no further
tasks will be started for that workflow instance.

### `output`

You can obtain the output (if any) of a completed workflow run using the output command and supplying the workflow run
id. Typically, this is useful for locating the files produced by a workflow, although the actual output generated depends on the workflow
specification and engine.

If the workflow declares outputs you may also obtain these using the command:

```shell
agc workflow output <workflow_run_id>
```

The following is an example of output from the "CramToBam" workflow run in a context using the Cromwell engine.

```shell
OUTPUT	id	aaba95e8-7512-48c3-9a61-1fd837ff6099
OUTPUT	outputs.CramToBamFlow.outputBai	s3://agc-123456789012-us-east-1/project/GATK/userid/mrschre4GqyMA/context/spotCtx/cromwell-execution/CramToBamFlow/aaba95e8-7512-48c3-9a61-1fd837ff6099/call-CramToBamTask/NA12878.bai
OUTPUT	outputs.CramToBamFlow.outputBam	s3://agc-123456789012-us-east-1/project/GATK/userid/mrschre4GqyMA/context/spotCtx/cromwell-execution/CramToBamFlow/aaba95e8-7512-48c3-9a61-1fd837ff6099/call-CramToBamTask/NA12878.bam
OUTPUT	outputs.CramToBamFlow.validation_report	s3://agc-123456789012-us-east-1/project/GATK/userid/mrschre4GqyMA/context/spotCtx/cromwell-execution/CramToBamFlow/aaba95e8-7512-48c3-9a61-1fd837ff6099/call-ValidateSamFile/NA12878.validation_report
```

## Cost

Your account will be charged based on actual resource usage including compute time, storage, data transfer charges etc.
The resources used will depend on the resources requested in your workflow definition as interpreted by the workflow engine
according the resources made available in the context in which the workflow is run. If a spot context is used then the costs
of the spot instances will be determined by the rules governing spot instance charges.

### Tags

Resources used by Amazon Genomics CLI are tagged including the username, project name and the context name. Currently, tagging is *not*
possible at the level of an individual workflow.
