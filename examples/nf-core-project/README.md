## Introduction

This project uses workflows from the [nf-core](https://nf-co.re/) repository. For licensing terms, credits and citation of individual
workflows please refer to the following links:

* https://github.com/nf-core/rnaseq
* https://github.com/nf-core/sarek
* https://github.com/nf-core/atacseq

## Example Workflow Runtimes

These times were recorded using the project's `bigMemCtx` context in `us-east-1` with a maximum of 256 vCPUs. Times include the time from workflow submission
to completion and for the AWS Batch service to provision compute resources (cold start time). Not included is the time required for the Nextflow headnode startup.
The results are indicative only and runtimes will vary based on the resources allocated by AWS Batch as well as any automated retries
of tasks due to task failures. Caching of tasks was not used in this evaluation.

Workflow inputs were those defined in the respective `*inputs.json` files


| Workflow | Runtime (minutes) |
| -------- | ----------------- |
| atacseq | 139 |
| rnaseq | 458 |
| sarek | 359 |