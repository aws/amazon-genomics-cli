---
title: "Overview"
linkTitle: "Overview"
weight: 1
description: >
  Amazon Genomics CLI Overview.
---

{{% alert title="Attention" color="warning" %}}
**The Amazon Genomics CLI project has entered its End Of Life (EOL) phase**. The code is no longer actively maintained and the **Github repository will be archived on May 31 2024**. During this time, we encourage customers to migrate to [AWS HealthOmics](https://aws.amazon.com/healthomics/) to run their genomics workflows on AWS, or [reach out to their AWS account team](https://aws.amazon.com/contact-us/?nc2=h_header) for alternative solutions. While the source code of AGC will still be available after the EOL date, we will not make any updates inclusive of addressing issues or accepting Pull Requests.
{{% /alert %}}

## What is the Amazon Genomics CLI?

Amazon Genomics CLI is an open source tool for genomics and life science customers that simplifies and automates the 
deployment of cloud infrastructure, providing you with an easy-to-use command line interface to quickly setup and run 
genomics workflows on Amazon Web Services (AWS) specified by languages like CWL, Nextflow, Snakemake, and WDL. By removing the heavy lifting from
setting up and running genomics workflows in the cloud, software developers and researchers can automatically provision, 
configure and scale cloud resources to enable faster and more cost-effective population-level genetics studies, drug 
discovery cycles, and more.

## Why do I want it?

Amazon Genomics CLI is targeted at bioinformaticians and genomics analysts who are not experts in cloud infrastructure
and best practices. If you want to take advantage of cloud computing to run your workflows at scale, but you don't want
to become an expert in high performance batch computing and distributed systems then Amazon Genomics CLI is probably
for you.

* **What is it good for?**: Abstracting the infrastructure needed to run workflows from the running of the workflows by hiding all the complexity behind a familiar CLI interface. When you need to get your workflows running quickly and are happy to let the tool make the decisions about the best infrastructure according to AWS best practices.

* **What is it not good for?**: Situations where you want, or need, complete and fine-grained control over how your workflows are run in the cloud. Where specifying *exactly* how they run, and what infrastructure is used, is just as important as running them.

* **What is it *not yet* good for?**: [Let us know](https://github.com/aws/amazon-genomics-cli/issues/new/choose), we'd love to hear your suggestions on how we can make the tool work for you.

## Where should I go next?

Take a look at the following pages to help you start running your workflows quickly:

* [Getting Started]({{< ref "/docs/Getting started/" >}}): Get started with Amazon Genomics CLI
* [Tutorials]( {{< ref "/docs/Tutorials/" >}} ): Tutorials to show you the ropes.

