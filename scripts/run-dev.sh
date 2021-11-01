#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

set -eo pipefail
trap cleanup SIGQUIT
trap cleanup SIGTSTP

TEMP_STORAGE=$(mktemp -d 2>/dev/null || mktemp -d -t 'agc')

cleanup () {
  echo ""
  echo "Cleaning up and exiting"
  rm -rf "$HOME/.agc/"
  mv "$TEMP_STORAGE/.agc" "$HOME"
  rm -rf "$TEMP_STORAGE"
  exit 0
}

get_latest_version () {
  if [[ -z "$1" || -z "$2" || -z "$3" ]]; then
    >&2 echo "Repository, account, and region must be specified, but got '$1', '$2', and '$3'"
    exit 1
  fi

  imagesJson=$(aws ecr list-images --repository-name "$1" --registry-id "$2" --region "$3" --filter tagStatus=TAGGED --output json)
  if [[ $? -ne 0 ]]; then
    >&2 echo "Could not list images for repository '$1', account '$2', and region '$3'"
    exit 1
  fi

  latestImage=$(echo "$imagesJson" | jq -r '.imageIds[].imageTag' | sort | tail -1)
  echo "${latestImage}"
}

run_agc_cli () {
  trap run_agc_cli SIGINT
  while true; do
    read -rea command -p "$ agc "
    history -s "${command[@]}"
    "${SCRIPT_DIR}"/../packages/cli/bin/local/agc "${command[@]}" || true
  done
}

WES_REPOSITORY="agc-wes-adapter-cromwell"
CROMWELL_REPOSITORY="cromwell"
NEXTFLOW_REPOSITORY="nextflow"

ECR_WES_ACCOUNT_ID="555741984805"
ECR_WES_REGION="us-east-1"
export ECR_WES_ACCOUNT_ID ECR_WES_REGION

ECR_CROMWELL_ACCOUNT_ID="555741984805"
ECR_CROMWELL_REGION="us-east-1"
export ECR_CROMWELL_ACCOUNT_ID ECR_CROMWELL_REGION

ECR_NEXTFLOW_ACCOUNT_ID="555741984805"
ECR_NEXTFLOW_REGION="us-east-1"
export ECR_NEXTFLOW_ACCOUNT_ID ECR_NEXTFLOW_REGION

echo "Getting latest tag versions"
ECR_WES_TAG=$(get_latest_version "$WES_REPOSITORY" "$ECR_WES_ACCOUNT_ID" "$ECR_WES_REGION")
ECR_CROMWELL_TAG=$(get_latest_version "$CROMWELL_REPOSITORY" "$ECR_CROMWELL_ACCOUNT_ID" "$ECR_CROMWELL_REGION")
ECR_NEXTFLOW_TAG=$(get_latest_version "$NEXTFLOW_REPOSITORY" "$ECR_NEXTFLOW_ACCOUNT_ID" "$ECR_NEXTFLOW_REGION")

export ECR_WES_TAG ECR_CROMWELL_TAG ECR_NEXTFLOW_TAG

export LOCAL_CORE_SERVICE=true

echo "The following environment variables have been set"
env | grep ECR | sort

echo "Setting up CDK"
mkdir -p "$HOME/.agc/"
mv "$HOME/.agc" "$TEMP_STORAGE"
mkdir -p "$HOME/.agc/"
echo -e "user:\n    email: $USER@amazon.com" > "$HOME/.agc/config.yaml"
ln -sfn "${SCRIPT_DIR}/../packages/cdk/" "$HOME/.agc/"

echo "Running the agc CLI"
echo "  ^C will interrupt the agc command"
echo "  To exit this script, use ^\ or ^Z"
echo ""
run_agc_cli
