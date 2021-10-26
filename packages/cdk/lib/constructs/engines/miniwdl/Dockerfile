# use the upstream nextflow container as a base image
ARG NEXTFLOW_VERSION=20.10.0
FROM public.ecr.aws/amazonlinux/amazonlinux:latest as build

RUN amazon-linux-extras install -y python3.8
RUN yum update -y \
 && yum install -y \
    unzip \
    jq \
 && yum clean -y all
RUN rm -rf /var/cache/yum

# install awscli v2
RUN curl -s "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip" \
 && unzip -q /tmp/awscliv2.zip -d /tmp \
 && /tmp/aws/install -b /usr/bin \
 && rm -rf /tmp/aws*


RUN python3.8 -m venv /venv
ENV PATH="/venv/bin:$PATH"

RUN pip install --upgrade pip
RUN pip install miniwdl-aws

# install a custom entrypoint script that handles being run within an AWS Batch Job
COPY miniwdl.aws.sh /opt/bin/miniwdl.aws.sh
RUN chmod +x /opt/bin/miniwdl.aws.sh

WORKDIR /mnt/efs
ENTRYPOINT ["/opt/bin/miniwdl.aws.sh"]
