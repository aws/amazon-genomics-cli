FROM public.ecr.aws/amazonlinux/amazonlinux:2 AS final

# COPY THIRD-PARTY /opt/
COPY LICENSE /opt/

RUN yum update -y \
    && yum install -y \
    curl \
    hostname \
    "java-11-amazon-corretto-headless(x86-64)" \
    unzip \
    jq \
    && yum clean -y all \
    && rm -rf /var/cache/yum

# install awscli v2
RUN curl -s "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip" \
    && unzip -q /tmp/awscliv2.zip -d /tmp \
    && /tmp/aws/install -b /usr/bin \
    && rm -rf /tmp/aws*

##### MODIFY #######
## In this area install your new engine into the container as well as any requirements for that engine.
## Dockerfile documentation is found here: https://docs.docker.com/engine/reference/builder/

# Add rabbitmq repository
ADD rabbitmq.repo /etc/yum.repos.d/rabbitmq.repo

# Sadly pre-importing keys doesn't seem to save any time when we use yum later, so don't so it.

# Install deps
RUN curl -fsSL https://rpm.nodesource.com/setup_16.x | bash - \
    && yum update -y \
    && yum install -y \
    python3 \
    rabbitmq-server \
    nodejs \
    git \
    && yum clean -y all \
    && rm -rf /var/cache/yum

# Install concurrently, for running all our servers in one session
RUN npm install -g concurrently@7.0.0

# Install Toil
COPY THIRD-PARTY /opt/

ARG TOIL_VERSION="d831f74e918c4a01e961e3b45504a92d1827b8b3"
RUN python3 -m pip install git+https://github.com/DataBiosphere/toil.git@${TOIL_VERSION}#egg=toil[aws,cwl,server]

# copy the entrypoint script to the image
COPY toil.aws.sh /opt/bin/toil.aws.sh
RUN chmod +x /opt/bin/toil.aws.sh

EXPOSE 8000

#### END MODIFY ######

WORKDIR /opt/work
ENTRYPOINT ["/opt/bin/toil.aws.sh"]

