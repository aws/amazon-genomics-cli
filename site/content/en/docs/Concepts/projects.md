---
title: "Projects"
date: 2021-08-31T17:27:06-04:00
draft: false
weight: 10
description: >
    A project defines the contexts, engines, data and workflows that make up a genomics analysis
---
An Amazon Genomics CLI project defines the [projects]( {{< relref "projects" >}} ), [contexts]( {{< relref "contexts" >}} ), [data]( {{< relref "data" >}} ) and [workflows]( {{< relref "workflows" >}} ) that make up a genomics analysis. Each project is defined
in a project file named `agc-project.yaml`.

## Project File Location

To find the project definition, Amazon Genomics CLI will look for a file named `agc-project.yaml` in the current working directory. If
the file is not found, Amazon Genomics CLI will traverse up the file hierarchy until the file is found or until the root of the file
system is reached. If no project definition can be found an error will be reported. All Amazon Genomics CLI commands operate on the project identified by the above process.

Consider the example directory structure below:

```
/
├── baa/
│   ├── a/
│   └── agc-project.yaml
├── foo/
└── foz/
    └── a/
        ├── agc-project.yaml
        └── b/
            └── c/
                └── agc-project.yaml
```

* If the current working directory is `/baa` or `/baa/a` then `/baa/agc-project.yaml` will be used for definitions,
* If the current working directory is `/foo` an error will be reported as no project file is found before the root,
* If the current working directory is `/foz` an error will be reported as no project file is found before the root,
* If the current working directory is `/foz/a` or `/foz/a/b` then `/foz/a/agc-project.yaml` will be used for definitions.
* If the current working directory is `/foz/a/b/c` then `/foz/a/b/c/agc-project.yaml` will be used for definitions.

### Relative Locations
The location of resources declared in a project file are resolved relative to the location of the project file *unless*
they are declared using an absolute path. If the project file in `/baa` declared that 
there was a workflow definition in `a/b/` then Amazon Genomics CLI will search for that definition in `/baa/a/b/`. 

## Project File Structure

A minimal project file can be generated using the `agc project init myProject --workflow-type nextflow`. Using `myProject` as a project name and workflow type `nextflow` will result in the following:

```yaml
name: myProject
schemaVersion: 1
contexts:
  ctx1:
    engines:
      - type: nextflow
        engine: nextflow
```

This is fully usable project called "myProject" with a single context named "ctx1". At this point "ctx1" can be [deployed]( {{< relref "contexts#deploy" >}} )
however, there are currently no workflows defined.

### `name`

A string that identifies the project

### `schemaVersion`

An integer defining the schema version. Version numbers will be incremented when changes are made to the project schema
that are not backward compatible.

### `contexts`

A map of context names to context definitions. Each context in the project must have a unique name. The [contexts]( {{< relref "contexts" >}} )
documentation provides more details.

### `workflows`

A map of workflow names to workflow definitions. Workflow names must be unique in a project. The [workflows]( {{< relref "workflows" >}} )
documentation provides more details.

### `data`

An array of data sources that the contexts of the project have access to. For example:

```yaml
data:
  - location: s3://gatk-test-data
    readOnly: true
  - location: s3://broad-references
    readOnly: true
  - location: s3://1000genomes-dragen-3.7.6
    readOnly: true
```

## Commands

A full reference of project commands are available [here]( {{< relref "../Reference/agc_project" >}} )

### `init`

The `agc project init <project-name> --workflow-type <worklow-type>` command can be used to initialize a minimal `agc-project.yaml` file in the current
directory. Alternatively project yaml files can be created with any text editor.

### `describe`

The `agc project describe <project-name>` command will provide basic metadata about the 'local' project file. See 
[above](#project-file-location) for details on how project files are located.

### `validate`

Using `agc project validate` you can quickly identify any syntax errors in your local project file.

## Versioning and Sharing

We recommend placing a project under source version control using a tool like [Git](https://git-scm.com). The folder containing the `agc-project.yaml`
file is a natural location for the root of a Git repository. Workflows relating to the project would naturally be located 
in sub-folders of the same repository allowing those to be versioned as well. Alternatively, more advanced Git users may
consider storing workflows as a Git [sub-module](https://git-scm.com/book/en/v2/Git-Tools-Submodules) allowing them to 
be independent of the project and reused among projects.

Projects and associated workflows can then be shared by "pushing" the project's Git repository to a website such
as [GitHub](https://github.com), [GitLab](https://gitlab.com), or [BitBucket](https://bitbucket.com) or hosted on a 
private Git Server like [AWS Code Commit](https://docs.aws.amazon.com/codecommit/latest/userguide/index.html). To facilitate sharing you should
ensure that any file paths in your definitions are relative to the project and not absolute. You will also need to make
sure that data locations are appropriately shared.

## Costs

A project itself doesn't have infrastructure. It is not deployed and therefore has no direct costs. If the contexts defined
by an infrastructure are deployed or the workflows run then those *will* incur costs. 

### Tags

The project `name` will be [tagged]( {{< relref "namespaces#tags" >}} )
on any deployed contexts or workflows defined in this project allowing costs to be aggregated to the project level.

## Technical Details

A project is purely a [YAML](https://en.wikipedia.org/wiki/YAML) definition. The values in the `agc-project.yaml` file are used by CDK when Amazon Genomics CLI deploys contexts
and when Amazon Genomics CLI runs workflows. The project itself has no direct infrastructure. The project `name` is used to help namespace
context infrastructure.
