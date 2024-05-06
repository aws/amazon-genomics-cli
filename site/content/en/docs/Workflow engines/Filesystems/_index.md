---
title: "Filesystems"
linkTitle: "FileSystems"
weight: 10
date: 2022-03-11T14:31:11-04:00
draft: false
description: >
    Workflow Filesystems
---

{{% alert title="Attention" color="warning" %}}
**The Amazon Genomics CLI project has entered its End Of Life (EOL) phase**. The code is no longer actively maintained and the **Github repository will be archived on May 31 2024**. During this time, we encourage customers to migrate to [AWS HealthOmics](https://aws.amazon.com/healthomics/) to run their genomics workflows on AWS, or [reach out to their AWS account team](https://aws.amazon.com/contact-us/?nc2=h_header) for alternative solutions. While the source code of AGC will still be available after the EOL date, we will not make any updates inclusive of addressing issues or accepting Pull Requests.
{{% /alert %}}

The tasks in a workflow require a common filesystem or scratch space where the outputs of tasks can be written so they 
are available to the inputs of dependent tasks in the same workflow. The following pages provide details on the engine 
filesystems that can be deployed by Amazon Genomics CLI.

