## Nextflow AWS Mirror

An AWS-friendly mirror of Nextflow.

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
