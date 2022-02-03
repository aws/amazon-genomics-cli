#!/bin/bash

set -e
set -x

OS=$(uname -r)
export OS
BASEDIR=$(dirname "$0")
INITIAL_EBS_SIZE="${1:-200}"

echo OS = "$OS"
echo BASEDIR = "$BASEDIR"
echo INITIAL_EBS_SIZE = "$INITIAL_EBS_SIZE"
echo ARTIFACTS_NAMESPACE = "$ARTIFACTS_NAMESPACE"


# Expected environment variables
#   ARTIFACT_S3_ROOT_URL (obtained from SSM parameter store)
#   WORKFLOW_ORCHESTRATOR (OPTIONAL)

printenv

function ecs() {
    
    if [[ $OS =~ "amzn1" ]]; then
        # Amazon Linux 1 uses upstart for init
        case $1 in
            disable)
                echo "stopping ecs service"
                stop ecs
                echo "stopping docker service"
                service docker stop
                ;;
            enable)
                echo "starting docker service"
                service docker start
                echo "starting ecs service"
                start ecs
                ;;
        esac
    elif [[ $OS =~ "amzn2" ]]; then
        # Amazon Linux 2 uses systemd for init
        case $1 in
            disable)
                echo "stopping ecs service"
                systemctl stop ecs
                echo "stopping docker service"
                systemctl stop docker
                ;;
            enable)
                echo "starting docker service"
                systemctl start docker
                echo "enabling ecs service"
                systemctl enable --now --no-block ecs  # see: https://github.com/aws/amazon-ecs-agent/issues/1707
                ;;
        esac
    else
        echo "unsupported os: $os"
        exit 100
    fi
}

function getArtifactRoot(){
  url=$(\
      aws ssm get-parameter \
          --name /"${ARTIFACTS_NAMESPACE}"/_common/installed-artifacts/s3-root-url \
          --query 'Parameter.Value' \
          --output text \
  )
  return "$url"
}

function errorOrInt() {
  echo "WARNING - received a $1 signal. Will attempt to restart ECS agent but this instance may not be correctly provisioned"
  env
  ecs enable
}

# make sure that docker and ecs are running on script exit to avoid
# zombie instances
trap "ecs enable" EXIT
trap "errorOrInt INT" INT
trap "errorOrInt ERR" ERR

set +e
ecs disable
set -e

ARTIFACT_S3_ROOT_URL=getArtifactRoot
if [[ -v ${ARTIFACT_S3_ROOT_URL} ]]; then
    echo "ARTIFACT_S3_ROOT_URL not found, trying again" && sleep 5 && ARTIFACT_S3_ROOT_URL=getArtifactRoot
fi
echo "ARTIFACT_S3_ROOT_URL = $ARTIFACT_S3_ROOT_URL"


# retrieve and install amazon-ebs-autoscale
echo "WORKFLOW_ORCHESTRATOR = $WORKFLOW_ORCHESTRATOR"
if [ "$WORKFLOW_ORCHESTRATOR" != "miniwdl" ]; then
  echo "obtaining amazon-ebs-autoscale artifacts"
  cd /opt
  sh "$BASEDIR"/get-amazon-ebs-autoscale.sh \
      --install-version dist_release \
      --artifact-root-url "$ARTIFACT_S3_ROOT_URL" \
      --file-system btrfs \
      --initial-size "$INITIAL_EBS_SIZE"
  echo "amazon-ebs-autoscale artifacts installed with return code = $?"
fi

# common provisioning for all workflow orchestrators
cd /opt
echo "installing ecs-additions-common"
sh "$BASEDIR"/ecs-additions-common.sh
echo "installed ecs-additions-common"

# workflow specific provisioning if needed
if [[ $WORKFLOW_ORCHESTRATOR ]]; then
    if [ -f "$BASEDIR/ecs-additions-$WORKFLOW_ORCHESTRATOR.sh" ]; then
        echo "installing orchestrator specific provisioning $BASEDIR/ecs-additions-$WORKFLOW_ORCHESTRATOR.sh"
        sh "$BASEDIR"/ecs-additions-"$WORKFLOW_ORCHESTRATOR".sh
        echo "orchestrator specific provisioning complete with return code = $?"
    fi
fi
