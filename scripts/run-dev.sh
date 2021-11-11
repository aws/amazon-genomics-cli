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

run_agc_cli () {
  trap run_agc_cli SIGINT
  while true; do
    read -rea command -p "$ agc "
    history -s "${command[@]}"
    "${SCRIPT_DIR}"/../packages/cli/bin/local/agc "${command[@]}" || true
  done
}

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
