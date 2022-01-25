#! /bin/bash

rm -rf ./package && mkdir ./package
pip3 install -r requirements.txt --target ./package && cd ./package && zip -r ../wes_adapter.zip .
cd ..
zip -gr ./wes_adapter.zip ./rest_api
zip -gr ./wes_adapter.zip ./amazon_genomics
zip -g ./wes_adapter.zip ./index.py
