#!/bin/bash

# Toil is a WES server and so it gets this custom entrypoint script

DEFAULT_AWS_CLI_PATH=/opt/aws-cli/bin/aws
AWS_CLI_PATH=${JOB_AWS_CLI_PATH:-$DEFAULT_AWS_CLI_PATH}

echo "=== ENVIRONMENT ==="
printenv

echo "=== START SERVER ==="

# We expect some AGC info in the environment: JOB_QUEUE_ARN
# These come from packages/cdk/lib/env/context-app-parameters.ts
# If we need more we'll need to add them in the Toil engine construct, or maybe stop passing getEngineContainer() down as a parameter.
# We assume whatever role the batch jobs get when they go in the queue is the right role for them.
AWS_REGION=$(echo ${JOB_QUEUE_ARN} | cut -f4 -d':')
set -x

export TOIL_WES_BROKER_URL="amqp://guest:guest@localhost:5672//"
export TOIL_WES_JOB_STORE_TYPE="aws"

concurrently -n rabbitmq,celery,toil \
    "rabbitmq-server" \
    "celery --broker=${TOIL_WES_BROKER_URL} -A toil.server.celery_app worker --loglevel=INFO" \
    "toil server --debug --host=0.0.0.0 --port=8000 --opt=--batchSystem=aws_batch '--opt=--awsBatchQueue=${JOB_QUEUE_ARN}' '--opt=--awsBatchRegion=${AWS_REGION}'"


