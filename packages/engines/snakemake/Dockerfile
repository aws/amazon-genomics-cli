ARG SNAKEMAKE_VERSION
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

# install conda 
RUN curl -L https://repo.continuum.io/miniconda/Miniconda3-latest-Linux-x86_64.sh > miniconda.sh && \
    bash miniconda.sh -b -p /opt/conda && \
    rm miniconda.sh
ENV PATH /opt/conda/bin:${PATH}

ENV JAVA_HOME /usr/lib/jvm/jre-openjdk/

RUN conda install -y -c conda-forge mamba && \
    mamba create -q -y -c conda-forge -c bioconda -n snakemake snakemake=${SNAKEMAKE_VERSION} && \
    source activate snakemake

COPY snakemake.aws.sh /opt/bin/snakemake.aws.sh
RUN chmod +x /opt/bin/snakemake.aws.sh

ENV PATH /opt/conda/envs/snakemake/bin:${PATH}

WORKDIR /opt/work
ENTRYPOINT ["/opt/bin/snakemake.aws.sh"]