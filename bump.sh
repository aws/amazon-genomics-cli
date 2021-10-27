#!/bin/bash

#--------------------------------------------------------------------------------------------------
#
# This script is intended to be used to bump up the version of the AGC modules for major, minor, and patch releases
# --------------------------------------------------------------------------------------------------

set -euo pipefail
releaseAs=${1:-}
if [ -z "${releaseAs}" ] || [ "${releaseAs}" == "major" ]; then
  echo "usage: ./bump.sh <version> or minor or patch. For major release, please specify the exact version number expected after bumping up the version."
  exit 1
fi

npx standard-version --release-as "$releaseAs"
