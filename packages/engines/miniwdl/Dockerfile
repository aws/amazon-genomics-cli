# use the upstream miniwdl container as a base image
ARG MINIWDL_VERSION=v0.1.11
FROM 337532070941.dkr.ecr.eu-central-1.amazonaws.com/miniwdl-aws:latest

RUN yum update -y \
 && yum install -y \
    unzip \
    jq \
 && yum clean -y all
RUN rm -rf /var/cache/yum

COPY THIRD-PARTY /opt/
COPY LICENSE /opt/
COPY miniwdl.aws.sh /opt/bin/miniwdl.aws.sh
RUN chmod +x /opt/bin/miniwdl.aws.sh

WORKDIR /opt/work
ENTRYPOINT ["/opt/bin/miniwdl.aws.sh"]
