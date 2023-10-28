#!/usr/bin/env bash

set -eo pipefail

USER_BIN_DIR="$HOME/bin"
BASE_DIR="$HOME/.agc"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}/../../" )" &> /dev/null && pwd )"

uninstall_cli () {
    rm -f $USER_BIN_DIR/agc
}

uninstall_cdk () {
    rm -rf $BASE_DIR
}

uninstall_cdk && uninstall_cli && echo "AGC has been uninstalled."
