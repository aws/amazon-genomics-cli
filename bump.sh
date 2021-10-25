#!/usr/bin/env bash

set -eo pipefail

releaseAs="${1:-minor}"
if [ "$releaseAs" != "minor" ] && [ "$releaseAs" != "patch" ]; then
  echo "Only 'minor' and 'patch' version bumps are allowed. Major version bumps will require more intention."
  exit 1
fi


echo "Releasing a $releaseAs version bump..."
npx standard-version --release-as "$releaseAs"