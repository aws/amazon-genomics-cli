#!/usr/bin/env bash

set -eo pipefail

USER_BIN_DIR="$HOME/bin"
BASE_DIR="$HOME/.agc"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

selectArch () {
    if [[ $(arch) == "arm64" ]]; then
        eval $1="arm64"
    else
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            eval $1="amd64"
        fi
    fi
}

selectCliFile () {
    selectArch archKind
    fileName="agc"

    if [[ -z "$archKind" ]]; then
        eval $fileName="agc"
    else
        eval $fileName="$fileName-$archKind"
    fi

    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        eval $1=$fileName
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        eval $1=$fileName
    else
        echo "You are running on ${OSTYPE}. AGC does not yet support this platform."
        echo "Please try macOS or a Debian based OS."
    fi
}

install_cli () {
    selectCliFile cliFile
    if [[ -z "$cliFile" ]]; then
        exit 1
    fi

    mkdir -p $USER_BIN_DIR
    cp "$SCRIPT_DIR/$cliFile" "$USER_BIN_DIR/agc"

    echo "Please modify your \$PATH variable to include \$HOME/bin directory"
    echo "This can be achieved by running: \"export PATH=\$HOME/bin:\$PATH\""
    echo "Please append the command above to shell profile to have agc available within every shell instance." 
}

install_cdk () {
    mkdir -p "$BASE_DIR/cdk"
    cp "$SCRIPT_DIR/cdk.tgz" "$BASE_DIR/cdk"
    (cd "$BASE_DIR/cdk" && tar -xzf ./cdk.tgz --strip-components=1 && npm ci --silent)
}

install_wes () {
    mkdir -p "$BASE_DIR/wes"
    cp "$SCRIPT_DIR/wes/wes_adapter.zip" "$BASE_DIR/wes"
}

install_cli && install_cdk && install_wes && echo "Installation complete. Once \$PATH variable has been adjusted, run 'agc --help' to get started!"
