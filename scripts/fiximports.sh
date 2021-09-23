#!/usr/bin/env bash

set -eo pipefail

find . -type f -name "*.go" -not -path "*/vendor*" -exec sed -i '
  /^import/,/)/ {
    /^$/ d
  }
' "{}" +
