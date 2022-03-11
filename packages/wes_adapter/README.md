# Amazon Genomics Workflow Execution Service (WES) Adapter

A basic GA4GH compliant Workflow Execution Service that enables use of adaptors for integrating with workflow execution engines.
The code in this package is capturing the Lambda Function source code that's deployed as WES adapter endpoint for purposes of AGC


# Open source usage
This package contains the source code that will be built and packaged before uploading into your AWS account as lambda code; this upload will contain the open source project dependencies shown in the requirement.txt file

# Local Run

This code could be run locally

## Step 1

Run `make init`, this will create a venv and install the packages required for WES Lambda to run

## Step 2

Override the required environment variables at `./start-local-server.sh` to point the service to access correct AWS resources

```bash
export ENGINE_NAME= # nextflow, snakemake, miniwdl or cromwell
export JOB_QUEUE= 
export JOB_DEFINITION= 
export ENGINE_LOG_GROUP=
```

## Step 3

Execute `./start-local-server.sh` and navigate your browser to http://localhost:80/ga4gh/wes/v1/ui/

# Make release

To run `make release`, you need to have python3.9 installed. 
There are several ways to install python3.9, the easiest way is to download installer from
https://www.python.org/downloads/.

## Mac OS
Mac users can install python with homebrew by running `brew install python@3.9`.

## Linux
Linux's users can install python with apt-get by running `apt-get install python3.9`.
