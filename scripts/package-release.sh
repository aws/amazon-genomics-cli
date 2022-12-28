#!/usr/bin/env bash

set -eo pipefail

RELEASE_DIR="dist/amazon-genomics-cli"

if ! command -v jq &> /dev/null
then
  echo "Missing required jq package"
  exit 1
fi

mkdir -p ${RELEASE_DIR}
cp ./{LICENSE,THIRD-PARTY,CHANGELOG.md} ${RELEASE_DIR}
cp packages/cdk/cdk.tgz ${RELEASE_DIR}
mkdir -p ${RELEASE_DIR}/wes
cp packages/wes_adapter/dist/wes_adapter.zip ${RELEASE_DIR}/wes
cp -a scripts/cli/. ${RELEASE_DIR}
cp -a examples ${RELEASE_DIR}
cp -a extras ${RELEASE_DIR}
cp -a packages/cli/bin/local/. ${RELEASE_DIR}
version=$(jq .version -r < version.json)
commit="${CODEBUILD_RESOLVED_SOURCE_VERSION:-$(git rev-parse --verify HEAD)}"

cat > ${RELEASE_DIR}/build.json <<HERE
{
  "name": "amazon-genomics-cli",
  "version": "${version}",
  "commit": "${commit}"
}
HERE

cat > ${RELEASE_DIR}/nightly-build.json <<HERE
{
  "name": "amazon-genomics-cli",
  "version": "nightly-build",
  "commit": "${commit}"
}
HERE
