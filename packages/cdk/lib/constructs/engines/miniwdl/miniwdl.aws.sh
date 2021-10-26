#!/bin/bash
# $1    Nextflow project. Can be an S3 URI, or git repo name.
# $2..  Additional parameters passed on to the nextflow cli

# using nextflow needs the following locations/directories provided as
# environment variables to the container
#  * NF_LOGSDIR: where caching and logging data are stored
#  * NF_WORKDIR: where intermmediate results are stored

set -e  # fail on any error

DEFAULT_AWS_CLI_PATH=/opt/aws-cli/bin/aws
AWS_CLI_PATH=${JOB_AWS_CLI_PATH:-$DEFAULT_AWS_CLI_PATH}

echo "=== ENVIRONMENT ==="
printenv

echo "=== RUN COMMAND ==="
echo "$@"

MINIWDL_PROJECT=$1
shift
MINIWDL_PARAMS="$@"

# AWS Batch places multiple jobs on an instance
# To avoid file path clobbering use the JobID and JobAttempt
# to create a unique path. This is important if /opt/work
# is mapped to a filesystem external to the container
GUID="$AWS_BATCH_JOB_ID/$AWS_BATCH_JOB_ATTEMPT"

if [ "$GUID" = "/" ]; then
    GUID=`date | md5sum | cut -d " " -f 1`
fi

WORK_DIR=/mnt/efs/$GUID
mkdir -p $WORK_DIR
cd $WORK_DIR
export MINIWDL_CFG=$WORK_DIR/miniwdl.cfg

cat << 'EOF' > $MINIWDL_CFG
[scheduler]
container_backend = aws_batch_job
call_concurrency = 100

[file_io]
root = /mnt/efs

[task_runtime]
# Default policy to retry spot-terminated jobs (up to three total attempts)
defaults = {
        "docker": "ubuntu:20.04",
        "preemptible": 2
    }

[call_cache]
# Cache call outputs in EFS folder (valid so long as all referenced input & output files remain
# unmodified on EFS).
dir = /mnt/efs/miniwdl_run/_CACHE/call
get = true
put = true

[download_cache]
dir = /mnt/efs/miniwdl_run/_CACHE/download
get = false
put = false

[aws]
# Last-resort job timeout for AWS Batch to enforce (attemptDurationSeconds)
job_timeout = 864000
# Internal rate-limiting periods (seconds) for AWS Batch API requests
# (may need to be increased if many concurrent workflow runs are planned)
describe_period = 1
submit_period = 1
# "Burn-in" to increase submit_period early on in the workflow run. May be useful to avoid
# overloading EFS/EBS/Batch during the frenzy that tends to occur at workflow startup (e.g. S3
# downloads of many input files), without necessarily harming scheduling QoS in steady state.
# Given burn-in factors B and C, the effective submit_period at T seconds elapsed in the workflow:
#    submit_period' = submit_period*max(1, C-T/B)
# So for example setting B=36, C=100 ramps the submit_period down from 100X to 1X over an hour.
# Exception: doesn't apply when there are zero submitted jobs in non-terminal states.
submit_period_b = 0
submit_period_c = 100

EOF


# stage in session cache
# .nextflow directory holds all session information for the current and past runs.
# it should be `sync`'d with an s3 uri, so that runs from previous sessions can be
# resumed
echo "== Restoring Session Cache =="
# aws s3 sync --no-progress $NF_LOGSDIR/.nextflow .nextflow

function preserve_session() {
    # stage out session cache
    if [ -d .nextflow ]; then
        echo "== Preserving Session Cache =="
        aws s3 sync --no-progress .nextflow $NF_LOGSDIR/.nextflow
    fi

    # .nextflow.log file has more detailed logging from the workflow run and is
    # nominally unique per run.
    #
    # when run locally, .nextflow.logs are automatically rotated
    # when syncing to S3 uniquely identify logs by the batch GUID
    if [ -f .nextflow.log ]; then
        echo "== Preserving Session Log =="
        aws s3 cp --no-progress .nextflow.log $NF_LOGSDIR/.nextflow.log.${GUID/\//.}
    fi
}

function show_log() {
    echo "=== Nextflow Log ==="
    cat $WORK_DIR/workflow.log
}

function cleanup() {
    set +e
    wait $MINIWDL_PID
    set -e
    echo "=== Running Cleanup ==="

    show_log
    #preserve_session

    echo "=== Bye! ==="
}

function cancel() {
    # AWS Batch sends a SIGTERM to a container if its job is cancelled/terminated
    # forward this signal to Nextflow so that it can cancel any pending workflow jobs

    set +e  # ignore errors here
    echo "=== !! CANCELLING WORKFLOW !! ==="
    echo "stopping miniwdl pid: $MINIWDL_PID"
    kill -TERM "$MINIWDL_PID"
    echo "waiting .."
    wait $MINIWDL_PID
    echo "=== !! cancellation complete !! ==="
    set -e
}

trap "cancel; cleanup" TERM
trap "cleanup" EXIT

# stage workflow definition
if [[ "$MINIWDL_PROJECT" =~ ^s3://.* ]]; then
    echo "== Staging S3 Project =="
    aws s3 cp --recursive --no-progress $MINIWDL_PROJECT ./project
    find . -name '*.zip' -execdir unzip -o '{}' ';'
    MINIWDL_PROJECT=./project
    echo "MINIWDL_PROJECT: $MINIWDL_PROJECT"
    ls $MINIWDL_PROJECT

    MANIFEST_JSON=${MINIWDL_PROJECT}/MANIFEST.json
    if test -f "$MANIFEST_JSON"; then
      echo "cat $MANIFEST_JSON"
      cat $MANIFEST_JSON
      MINIWDL_PROJECT="$WORK_DIR/project/$(cat $MANIFEST_JSON | jq -r '.mainWorkflowURL')"
      MINIWDL_PARAMS="${MINIWDL_PARAMS} $(cat $MANIFEST_JSON | jq -r '.engineOptions // empty')"
      INPUT_JSON="$WORK_DIR/project/$(cat $MANIFEST_JSON | jq -r '.inputFileURLs[0]')"
      if test -f "$INPUT_JSON"; then
        echo "cat $INPUT_JSON"
        cat "$INPUT_JSON"
        MINIWDL_PARAMS="${MINIWDL_PARAMS} --input ${INPUT_JSON}"
      fi
    fi
fi

echo "== Running Workflow =="
MINIWDL_PARAMS="${MINIWDL_PARAMS} --dir ${WORK_DIR}/."
echo "miniwdl run ${MINIWDL_PROJECT} ${MINIWDL_PARAMS}"

miniwdl-run-s3upload $MINIWDL_PROJECT $MINIWDL_PARAMS &
# miniwdl-aws-submit $MINIWDL_PROJECT $MINIWDL_PARAMS &

MINIWDL_PID=$!
echo "nextflow pid: $MINIWDL_PID"
jobs
echo "waiting .."
wait $MINIWDL_PID
