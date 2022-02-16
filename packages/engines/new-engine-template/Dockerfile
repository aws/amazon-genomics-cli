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
    && yum clean -y all
RUN rm -rf /var/cache/yum

# install awscli v2
RUN curl -s "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip" \
    && unzip -q /tmp/awscliv2.zip -d /tmp \
    && /tmp/aws/install -b /usr/bin \
    && rm -rf /tmp/aws*

##### MODIFY #######
## In this area install your new engine into the container as well as any requirements for that engine.
## Dockerfile documentation is found here: https://docs.docker.com/engine/reference/builder/


# copy the entrypoint script to the image
COPY <engine-name>.aws.sh /opt/bin/<engine-name>.aws.sh
RUN chmod +x /opt/bin/<engine>.aws.sh

# set the path for the new engine
ENV PATH <enginepath>:${PATH}

#### END MODIFY ######

WORKDIR /opt/work
ENTRYPOINT ["/opt/bin/<engine>.aws.sh"]