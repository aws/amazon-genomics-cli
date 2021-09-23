#!/usr/bin/env bash

set -eo pipefail

releaseBucketName="${1:-healthai-public-assets-us-east-1}"

tagVersion="$(git describe --tags --abbrev=0)"
tagVersion="${tagVersion#[vV]}"
tagVersion="${tagVersion:-0.0.0}"
tagMajorVersion="${tagVersion%%.*}"
tagMinorVersion="${tagVersion#*.}"
tagMinorVersion="${tagMinorVersion%.*}"

releaseVersion=$(aws s3 ls "${releaseBucketName}/amazon-genomics-cli/" --recursive | tr -s ' ' | cut -d ' ' -f4- | sort | tail -n 1 | xargs -r dirname | xargs -r basename)
releaseVersion="${releaseVersion#[vV]}"
releaseVersion="${releaseVersion:-0.0.0}"
releaseMajorVersion="${releaseVersion%%.*}"
releaseMinorVersion="${releaseVersion#*.}"
releaseMinorVersion="${releaseMinorVersion%.*}"
releasePatchVersion="${releaseVersion##*.}"
nextPatchVersion="$((releasePatchVersion+1))"

if [[ "$tagMajorVersion" -gt "$releaseMajorVersion" ]]; then
    semVer="$tagVersion"
elif [[ "$tagMinorVersion" -gt "$releaseMinorVersion" ]]; then
    semVer="$tagVersion"
else
    semVer="$releaseMajorVersion.$releaseMinorVersion.$nextPatchVersion"
fi

echo "$semVer"