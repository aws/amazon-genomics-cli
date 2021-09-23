#!/usr/bin/env bash

set -eo pipefail

TARGET_FILE="./packages/cli/environment/environment.go"
WES_ECR_TAG_PLACEHOLDER="WES_ECR_TAG_PLACEHOLDER"
CROMWELL_ECR_TAG_PLACEHOLDER="CROMWELL_ECR_TAG_PLACEHOLDER"
NEXTFLOW_ECR_TAG_PLACEHOLDER="NEXTFLOW_ECR_TAG_PLACEHOLDER"

find_and_replace () {
  if [[ -z "$1" || -z "$2" || -z "$3" ]]; then
    echo "File, target, and value must be specified, but got '$1', '$2', '$3'"
    exit 1
  fi
  sed -i -e "s/$2/$3/g" "$1"
}

find_and_replace "$TARGET_FILE" "$WES_ECR_TAG_PLACEHOLDER" "$WES_ECR_TAG"
find_and_replace "$TARGET_FILE" "$CROMWELL_ECR_TAG_PLACEHOLDER" "$CROMWELL_ECR_TAG"
find_and_replace "$TARGET_FILE" "$NEXTFLOW_ECR_TAG_PLACEHOLDER" "$NEXTFLOW_ECR_TAG"
