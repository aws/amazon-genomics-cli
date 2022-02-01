#! /bin/bash

if ! hash python3.9; then
    echo "python3.9 is not installed"
    exit 1
fi

python3.9 -m venv ./venv

./venv/bin/pip3.9 install waitress
./venv/bin/pip3.9 install -r requirements.txt
