#!/bin/bash
# $1    miniwdl project. Can be an S3 URI, or URL.
# $2..  Additional parameters passed on to the miniwdl cli

# using miniwdl needs the following locations/directories provided as
# environment variables to the container
#  * MINIWDL_S3_OUTPUT_URI: S3 URI pointing to miniwdl output root. Job outputs will be written to subdirectories.

set -e  # fail on any error

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

function cleanup() {
    set +e
    wait $MINIWDL_PID
    set -e
    echo "=== Running Cleanup ==="

    echo "=== Bye! ==="
}

function cancel() {
    # AWS Batch sends a SIGTERM to a container if its job is cancelled/terminated
    # forward this signal to MiniWdl so that it can cancel any pending workflow jobs

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
    MINIWDL_PROJECT_DIRECTORY="."
    aws s3 cp --no-progress $MINIWDL_PROJECT "${MINIWDL_PROJECT_DIRECTORY}/"
    find $MINIWDL_PROJECT_DIRECTORY -name '*.zip' -execdir unzip -o '{}' ';'
    ls -l $MINIWDL_PROJECT_DIRECTORY

    MANIFEST_JSON=${MINIWDL_PROJECT_DIRECTORY}/MANIFEST.json
    if test -f "$MANIFEST_JSON"; then
      echo "cat $MANIFEST_JSON"
      cat $MANIFEST_JSON
      MINIWDL_PROJECT="$(cat $MANIFEST_JSON | jq -r '.mainWorkflowURL')"

      if [[ $MINIWDL_PROJECT != *"://"* ]] ; then
        MINIWDL_PROJECT="${MINIWDL_PROJECT_DIRECTORY}/${MINIWDL_PROJECT}"
      fi
      MINIWDL_PARAMS="${MINIWDL_PARAMS} $(cat $MANIFEST_JSON | jq -r '.engineOptions // empty')"
      INPUT_FILE="${MINIWDL_PROJECT_DIRECTORY}/$(cat $MANIFEST_JSON | jq -r '.inputFileURLs[0] // empty')"
      if [[ -n "$INPUT_FILE" ]] ; then
         INPUT_JSON="${MINIWDL_PROJECT_DIRECTORY}/${INPUT_FILE}"
      fi
    else
      WDL_FILES=( $MINIWDL_PROJECT_DIRECTORY/*.wdl )
      if [ "${#WDL_FILES[@]}" -eq 1 ]; then
          MINIWDL_PROJECT=${WDL_FILES[0]}
      else
        echo "Ambiguous WDL entrypoint. Please specify a mainWorkflowURL in a MANIFEST.json:
         https://aws.github.io/amazon-genomics-cli/docs/concepts/workflows/#multi-file-workflows"
        exit 1
      fi
      INPUT_FILES=( $MINIWDL_PROJECT_DIRECTORY/*input*.json )
      if [ "${#INPUT_FILES[@]}" -le 1 ]; then
          INPUT_JSON=${INPUT_FILES[0]}
      else
        echo "Ambiguous input files. Please provide no more than one *input*.json file."
        exit 1
      fi
    fi
    if test -f "$INPUT_JSON"; then
      echo "cat $INPUT_JSON"
      cat $INPUT_JSON
      MINIWDL_PARAMS="${MINIWDL_PARAMS} --input ${INPUT_JSON}"
    fi
fi

echo "== Running Workflow =="
MINIWDL_PARAMS="${MINIWDL_PARAMS}
  --dir ${WORK_DIR}/.
  --s3upload ${MINIWDL_S3_OUTPUT_URI}/${AWS_BATCH_JOB_ID}/"
echo "miniwdl run ${MINIWDL_PROJECT} ${MINIWDL_PARAMS}"

miniwdl-run-s3upload $MINIWDL_PROJECT $MINIWDL_PARAMS &

MINIWDL_PID=$!
echo "miniwdl pid: $MINIWDL_PID"
jobs
echo "waiting .."
wait $MINIWDL_PID
