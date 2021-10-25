#!/bin/bash

set -eo pipefail

#----------------------------------------------------------------------------------------------------------
#
# This script is intended to be used to bump up the version of the AGC modules for minor and patch releases
# ----------------------------------------------------------------------------------------------------------

releaseAs=${1:-minor}
if [ "${releaseAs}" != "minor" ] && [ "${releaseAs}" != "patch" ]; then
  echo "usage: ./bump.sh minor or patch. Major version bumps will require more intention."
  exit 1
fi

echo "Releasing a $releaseAs version bump..."
npx standard-version --release-as "$releaseAs"
