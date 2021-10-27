#!/bin/bash

#----------------------------------------------------------------------------------------------------------
#
# This script is intended to be used to bump up the version of the AGC modules for minor and patch releases
# ----------------------------------------------------------------------------------------------------------

releaseAs=${1:-}
if [ "${releaseAs}" != "minor" ] || [ "${releaseAs}" != "patch" ]; then
  echo "usage: ./bump.sh minor or patch. Defaulting to a minor release"
  releaseAs="${1:-minor}"
fi

npx standard-version --release-as "$releaseAs"
