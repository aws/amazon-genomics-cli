---
title: "Prerequisites"
date: 2021-09-07T13:42:19-04:00
draft: false
weight: 10
---

To run Amazon Genomics CLI the following prerequisites must be met:

* A computer with one of the following operating systems:
  * macOS 10.14+
  * Amazon Linux 2
  * Ubuntu 20.04
* Internet access
* An AWS Account
* An AWS role with sufficient access. To generate the minimum required policies for admins and users, please follow the instructions [here](https://github.com/aws/amazon-genomics-cli/tree/main/extras/agc-minimal-permissions).

Running Amazon Genomics CLI on Windows has not been tested, but it should run in WSL 2 with Ubuntu 20.04

## Prerequisite installation

### Ubuntu 20.04

* Install node.js

```
curl -fsSL https://deb.nodesource.com/setup_15.x | sudo -E bash -
sudo apt-get install -y nodejs
```

* Install and configure AWS CLI

```
sudo apt install awscli
`aws configure`
`# ... set access key ID, secret access key, and region`
```

### Amazon Linux 2 (e.g. on an EC2 instance)

* Install node

```
`curl -sL https://rpm.nodesource.com/setup_16.x | sudo -E bash -
sudo yum install -y nodejs`
```

* If you have not already done so, configure your AWS credentials and default region

```
aws configure
```

### MacOS

* Install [Homebrew](https://brew.sh/)

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

* Install node

```
`brew install node`
```

* Install and configure AWS CLI

```
`brew install awscli`
aws configure
# ... set access key ID, secret access key, and region
```
