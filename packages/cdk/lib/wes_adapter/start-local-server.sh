#! /bin/bash
export ENGINE_NAME=nextflow # cromwell / miniwdl
export JOB_QUEUE=
export JOB_DEFINITION=
export ENGINE_LOG_GROUP=

export AWS_DEFAULT_REGION=us-west-2
export AWS_REGION=us-west-2

echo "Starting local WES endpoint at http://localhost:80/ga4gh/wes/v1/ui/"
./venv/bin/python3 ./local-server.py
