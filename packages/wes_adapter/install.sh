#! /bin/bash

rm -rf ./dist && mkdir ./dist
python3.9 -m pip install -r requirements.txt --target ./dist && (cd ./dist && zip -r ./wes_adapter.zip .)
zip -gr ./dist/wes_adapter.zip ./rest_api ./amazon_genomics ./index.py
