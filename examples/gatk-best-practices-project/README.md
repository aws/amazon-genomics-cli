## Introduction

This project includes workflows based on [GATK Best Practices](https://gatk.broadinstitute.org/hc/en-us), developed by
the [Broad Institute](https://www.broadinstitute.org/). More information on how these workflows work is available in
the [GATK Workflows Github repository](https://github.com/gatk-workflows).

## Example Workflow Runtimes

These "wall" times were recorded using the project's `spotCtx` context in `us-east-1` with a maximum of 256 vCPUs. Times include the time from workflow submission
to completion and include any time needed for the AWS Batch service to provision compute (cold start time).
The results are indicative only and runtimes will vary based on the resources allocated by AWS Batch as well as any automated retries
of tasks due to network errors or spot instance interruptions. No call caching was used when producing these timings.

The inputs used to run the workflows are those specified in the `inputs.json` files of their respective workflows

| Workflow | Runtime (minutes) |
| -------- | ----------------- |
| seq-format-validation | 19 |
| paired-fastq-to-unmapped-bam | 12 |
| interleaved-fastq-to-paired-fastq | 4 |
| cram-to-bam | 18 |
| bam-to-unmapped-bams | 15 |
| gatk4-data-processing | 32 |
| gatk4-germline-snps-indels | 29 |
| gatk4-rnaseq-germline-snps-indels | 236 |
| gatk4-basic-joint-genotyping | 189 |
