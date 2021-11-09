#! /bin/bash

./venv/bin/pip install black
./venv/bin/python3 -m black ./amazon_genomics
./venv/bin/python3 -m black ./rest_api/controllers
