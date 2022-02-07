# Amazon Genomics CLI

[![Join the chat at https://gitter.im/aws/amazon-genomics-cli](https://badges.gitter.im/aws/amazon-genomics-cli.svg)](https://gitter.im/aws/amazon-genomics-cli?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Overview

The Amazon Genomics CLI is a tool to simplify the processes of deploying the AWS infrastructure required to
run genomics workflows in the cloud, to submit those workflows to run, and to monitor the logs and outputs of those workflows.

## Quick Start

To get an introduction to Amazon Genomics CLI refer to the [Quick Start Guide](https://aws.github.io/amazon-genomics-cli/docs/getting-started/)
in our wiki.

## Further Reading

For full documentation, please refer to our [docs](https://aws.github.io/amazon-genomics-cli/docs/).

## Releases

All releases can be accessed on our releases [page](https://github.com/aws/amazon-genomics-cli/releases).

The latest nightly build can be accessed here: `s3://healthai-public-assets-us-east-1/amazon-genomics-cli/nightly-build/amazon-genomics-cli.zip`

## Development

To build from source you will need to ensure the following prerequisites are met.

### One-time setup

There are a few prerequisites you'll need to install on your machine before you can start developing.

Once you've installed all the dependencies listed here, run `make init` to install the rest.

#### Go

The Amazon Genomics CLI is written in Go.

To manage and install Go versions, we use [goenv](https://github.com/syndbg/goenv). Follow the installation
instructions [here](https://github.com/syndbg/goenv/blob/master/INSTALL.md).

Once goenv is installed, use it to install the version of Go required by the
Amazon Genomics CLI build process, so that it will be available when the build
process invokes goenv's `go` shim:

```bash
goenv install
```

You will need to do this step again whenever the required version of Go is
changed.

#### Node

Amazon Genomics CLI makes use of the AWS CDK to deploy infrastructure into an AWS account. Our CDK code is written in TypeScript.
You'll need Node to ensure the appropriate dependencies are installed at build time.

To manage and install Node versions, we use [nvm](https://github.com/nvm-sh/nvm).

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash
echo 'export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"' >> ~/.bashrc
echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm' >> ~/.bashrc
source ~/.bashrc
nvm install
```

Note: If you are using Zsh, replace `~/.bashrc` with `~/.zshrc`.

#### Sed (OSX)

OSX uses an outdated version of [sed](https://www.gnu.org/software/sed/manual/sed.html). If you are on a Mac, you will
need to use a newer version of sed to ensure script compatibility.

```bash
brew install gnu-sed
echo 'export PATH="$(brew --prefix gnu-sed)/libexec/gnubin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Note: If you are using Zsh, replace `~/.bashrc` with `~/.zshrc`.

#### One time setup

Once you've installed all the dependencies listed here, run `make init` to automatically install all remaining dependencies.


#### Make

We use `make` to build, test and deploy artifacts. To build and test issue the `make` command from the project root.

If you're experiencing build issues, try running `go clean --cache` in the project root to clean up your local go build cache. Then try to run `make init` then `make` again. This should ideally resolve it.

### Running Development Code

#### Option 1. Running with development script
To run against a development version of Amazon Genomics CLI, first build your relevant changes and then run `./scripts/run-dev.sh`. This will
set the required environment variables and then enter into an Amazon Genomics CLI command shell.

If you want to run from development code manually, ensure you have the following environment variables set.

```shell
export ECR_CROMWELL_ACCOUNT_ID=<some-value>
export ECR_CROMWELL_REGION=<some-value>
export ECR_CROMWELL_TAG=<some-value>
export ECR_NEXTFLOW_ACCOUNT_ID=<some-value>
export ECR_NEXTFLOW_REGION=<some-value>
export ECR_NEXTFLOW_TAG=<some-value>
export ECR_MINIWDL_ACCOUNT_ID=<some-value>
export ECR_MINIWDL_REGION=<some-value>
export ECR_MINIWDL_TAG=<some-value>
```

These environment variables point to the ECR account, region, and tags of the Cromwell, Nextflow, and MiniWDL engine respectively
that will be deployed for your contexts. They are written as Systems Manager Parameter Store variables when you activate
your Amazon Genomics CLI account region (`agc account activate`). The `./scripts/run-dev.sh` contains logic to determine the current
dev versions of the images which you would typically use. You may also use production images, the current values of which will
be written when you activate an account with the production version of Amazon Genomics CLI. If you have customized containers that you 
want to develop against you can specify these however you will need to make these available if you wish to make pull requests
with code that depends on them.

#### Option 2. Running with local release
Unlike running 'run-dev.sh' script, this option will build and install a new version of Amazon Genomics CLI, replacing 
the one installed. To run a release version of Amazon Genomics CLI from your local build, first build your changes and then run `make release`.
This will create a release bundle `dist/` at this package root directory. Run the `install.sh` script in the `dist` folder 
to install your local release version of Amazon Genomics CLI. After installing, you should be able to run `agc` on the terminal. 

### Building locally with CodeBuild

This package is buildable with AWS CodeBuild. You can use the AWS CodeBuild agent to run CodeBuild builds on a local
machine.

You only need to set up the build image the first time you run the agent, or when the image has changed. To set up the
build image, use the following commands:

```bash
git clone https://github.com/aws/aws-codebuild-docker-images.git
cd aws-codebuild-docker-images/ubuntu/standard/5.0
docker build -t aws/codebuild/standard:5.0 .
docker pull amazon/aws-codebuild-local:latest --disable-content-trust=false
```

Create an environment file (e.g. `env.txt`) with the appropriate entries depending on which image tags you want to use.

```shell
CROMWELL_ECR_TAG=2021-06-17T23-48-54Z
ECR_NEXTFLOW_TAG=2021-06-17T23-48-54Z
WES_ECR_TAG=2021-06-17T23-48-54Z
```

In the root directory for this package, download and run the CodeBuild build script:

```bash
wget https://raw.githubusercontent.com/aws/aws-codebuild-docker-images/master/local_builds/codebuild_build.sh
chmod +x codebuild_build.sh
./codebuild_build.sh -i aws/codebuild/standard:5.0 -a ./output -c -e env.txt
```

### Configuring docker image location

The default values for all variables are placeholders (e.g. 'WES_ECR_TAG_PLACEHOLDER'). It is replaces by the actual
value during a build process.

#### WES adapter for Cromwell

Local environment variables:

- `ECR_WES_ACCOUNT_ID`
- `ECR_WES_REGION`
- `ECR_WES_TAG`

The corresponding AWS Systems Manager Parameter Store property names:

- /agc/_common/wes/ecr-repo/account
- /agc/_common/wes/ecr-repo/region
- /agc/_common/wes/ecr-repo/tag

#### Cromwell engine

Local environment variables:

- `ECR_CROMWELL_ACCOUNT_ID`
- `ECR_CROMWELL_REGION`
- `ECR_CROMWELL_TAG`

The corresponding AWS Systems Manager Parameter Store property names:

- /agc/_common/cromwell/ecr-repo/account
- /agc/_common/cromwell/ecr-repo/region
- /agc/_common/cromwell/ecr-repo/tag

#### Nextflow engine

Local environment variables:

- `ECR_NEXTFLOW_ACCOUNT_ID`
- `ECR_NEXTFLOW_REGION`
- `ECR_NEXTFLOW_TAG`

The corresponding AWS Systems Manager Parameter Store property names:

- /agc/_common/nextflow/ecr-repo/account
- /agc/_common/nextflow/ecr-repo/region
- /agc/_common/nextflow/ecr-repo/tag

## Contributing

### Issues

See [Reporting Bugs/Feature Requests](CONTRIBUTING.md#reporting-bugsfeature-requests) for more information. For a list of
open bugs and feature requests, please refer to our [issues](https://github.com/aws/amazon-genomics-cli/issues?q=is%3Aopen+is%3Aissue) page.

### Pull Requests

See [Contributing via Pull Requests](CONTRIBUTING.md#contributing-via-pull-requests)

## Security

See [Security Issue Notification](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.

