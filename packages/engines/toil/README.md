## Toil AWS Mirror

A Toil mono-container WES server for use with Amazon Genomics CLI.

### Building the Container Manually (on AMD64)

Go to this directory and run:

```bash
docker build . -f Dockerfile -t toil-agc
```

### Running for Testing (on AMD64)

Having built the container, run:

```bash
docker run --name toil-agc-test -ti --rm -p "127.0.0.1:8000:8000" toil-agc
```

This will start the containerized server and make it available on port 8000 on the loopback interface. You can inspect the port mapping with:

```bash
docker port toil-agc-test
```

Then you can talk to it with e.g.:

```bash
curl -vvv "http://localhost:8000/ga4gh/wes/v1/service-info"
```

For debugging, you can get inside the container with:

```bash
docker exec -ti toil-agc-test /bin/bash
```

### Deploying (from AMD64)

To push this to an Amazon ECR repo, where AGC can get at it, you can do something like:

```bash
AWS_REGION=<your-deployment-region> # For example, us-west-2
AWS_ACCOUNT=<your-account-number> # For example, 123456789012
ECR_REPO=<your-ecr-repo> # For example, yourname/toil-agc. Needs to be created in the ECR console.
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com
docker build -t ${ECR_REPO} .
docker tag ${ECR_REPO}:latest ${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}:latest
docker push ${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}:latest
```

### Building and Deploying on ARM

If you are running on an ARM architecture system, `docker build` will try and build a Docker image for an ARM host. This will not work for two reasons. First, AGC uses AMD64 hosts and needs an AMD64 image. Second, the build will fail because the Erlang packages used by RabbitMQ are only available for AMD64. On ARM, the Erlang dependency will not be satisfiable:

```
#10 28.22 --> Processing Dependency: erlang >= 23.2 for package: rabbitmq-server-3.10.0-1.el7.noarch
#10 28.24 --> Finished Dependency Resolution
#10 28.26  You could try using --skip-broken to work around the problem
#10 28.26 Error: Package: rabbitmq-server-3.10.0-1.el7.noarch (rabbitmq_server)
#10 28.26            Requires: erlang >= 23.2
#10 28.27  You could try running: rpm -Va --nofiles --nodigest
------
executor failed running [/bin/sh -c curl -fsSL https://rpm.nodesource.com/setup_16.x | bash -     && yum update -y     && yum install -y     python3     rabbitmq-server     erlang     nodejs     git     && yum clean -y all     && rm -rf /var/cache/yum]: exit code: 1
```

To work around this, you will need to make sure that you have `docker buildx` installed and configured to be able to build AMD64 images, and then run:

```bash
AWS_REGION=<your-deployment-region> # For example, us-west-2
AWS_ACCOUNT=<your-account-number> # For example, 123456789012
ECR_REPO=<your-ecr-repo> # For example, yourname/toil-agc. Needs to be created in the ECR console.
docker buildx build --platform linux/amd64 --push --tag=${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}:latest -f Dockerfile .
```

This will build and push the image for the required architecture as a single operation.

### Using in AGC

To use a custom image in an AGC context, go to the directory for the project you want to use it with, set the environment variables, and deploy the context:

```bash
cd ../../../examples/demo-cwl-project/
export ECR_TOIL_ACCOUNT_ID="${AWS_ACCOUNT}"
export ECR_TOIL_TAG=latest
export ECR_TOIL_REPOSITORY="${ECR_REPO}"
export ECR_TOIL_REGION="${AWS_REGION}"
agc context deploy --context myContext
```
