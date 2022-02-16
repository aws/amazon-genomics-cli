## Snakemake AWS Mirror

An AWS-friendly mirror of Snakemake.

### Running locally with CodeBuild

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

In the root directory for this package, download and run the CodeBuild build script:

```bash
wget https://raw.githubusercontent.com/aws/aws-codebuild-docker-images/master/local_builds/codebuild_build.sh
chmod +x codebuild_build.sh
./codebuild_build.sh -i aws/codebuild/standard:5.0 -a ./output -c
```

### Running with docker

You can also build this package with docker locally.

1. Start by building the image with the following command:

```bash
docker build -t <imageName>:latest <location of Dockerfile>
```

2. Run with the following command:

```bash
docker run <imageName>
```

Notes:

The image uses the aws cli to copy files from s3 into the container. To do so the container requires aws credentials. One mechanims of providing these is to mount your local creds to the container. This
is done by running the `docker run` command as follows.

```
docker run -v ${HOME}/.aws/credentials:/root/.aws/credentials:ro <imageName>
```
