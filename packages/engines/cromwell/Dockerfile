# Running this container will require a number of environment variables with sensible values provided at run time. e.g.
# docker run -e ROOT_DIR=foo -e S3BUCKET="s3://foo" -e JOB_QUEUE_ARN="arn:aws::foo" -e AWS_REGION=us-east-1 cromwell:latest

FROM public.ecr.aws/amazonlinux/amazonlinux:2

COPY cromwell.jar cromwell.jar
COPY cromwell.conf cromwell.conf
COPY THIRD-PARTY /opt/
COPY LICENSE /opt/

RUN yum update -y
RUN yum install java-11-amazon-corretto-headless -y

ENTRYPOINT ["java", "-Dconfig.file=cromwell.conf", "-XX:MaxRAMPercentage=90.0", "-jar", "cromwell.jar", "server"]
