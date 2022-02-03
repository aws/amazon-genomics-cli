#!/bin/bash

#### **** For new engines, please only change MODIFY blocked areas below! ****
#### These areas are started with MODIFY and completed with END MODIFY

set -e  # fail on any error

DEFAULT_AWS_CLI_PATH=/opt/aws-cli/bin/aws
AWS_CLI_PATH=${JOB_AWS_CLI_PATH:-$DEFAULT_AWS_CLI_PATH}


echo "=== ENVIRONMENT ==="
printenv

# View paramaters to the script 
echo "=== RUN COMMAND ==="
echo "$@"

ENGINE_PROJECT=$1
shift
ENGINE_PARAMS="$@"

##### MODIFY #######
## set your engine name and the run command for the engine 

ENGINE_NAME=SNAKEMAKE
ENGINE_RUN_CMD=snakemake

#### END MODIFY ######
 
function handleManifest() {
   echo "cat $MANIFEST_JSON"
      cat $MANIFEST_JSON
      # Get correct url of project root location
      ENGINE_PROJECT="$(cat $MANIFEST_JSON | jq -r '.mainWorkflowURL')"
      if [[ $ENGINE_PROJECT != *"://"* ]] ; then
        ENGINE_PROJECT="${ENGINE_PROJECT_DIRECTORY}/${ENGINE_PROJECT}"
      fi
      ENGINE_OPTIONS="$(cat $MANIFEST_JSON | jq -r '.engineOptions')" 
      if [[ -n "$ENGINE_OPTIONS" ]] ; then
         ENGINE_PARAMS="${ENGINE_PARAMS} ${ENGINE_OPTIONS}"
      fi
      ###### MODIFY ###### 
        ## Add arguments to your run command if there's a manifest file
        ## In this area you can add any specialized paramaters the engine requires that come from the manifest file.
        ## To add to the beginning of the arguments, update PREPEND_ARGS.
        ## To add to the end of the arguments, update APPEND_ARGS.
        ## You can also use custom paramaters from the manifest.json file. For example:
          ## APPEND_ARGS="$(cat $MANIFEST_JSON | jq -r '.customArgument')" will pull the value attached to the key "customArgument"
          ## from the manifest.json file if the key exists, otherwise it will make it the value "".
        APPEND_ARGS="--aws-batch-tags AWS_BATCH_PARENT_JOB_ID=${AWS_BATCH_JOB_ID} --aws-batch-efs-project-path=snakemake/$GUID --latency-wait 60"
        PREPEND_ARGS=""
        ## To extend the arg strings you can do the following: APPEND_ARGS="${APPEND_ARGS} --moreArgumentsHere"
        ## After updating these you can expect your engine to be run as `ENGINE_RUN_CMD PREPEND_ARGS <agc_params> APPEND_ARGS`
        ## PREPEND_ARGS/APPEND_ARGS will only be added if they are not an empty string (the default above)
      ###### END MODIFY ######
      MODIFIED_PARAMS=false
      if [[ -n "$APPEND_ARGS" ]] ; then
        ENGINE_PARAMS="${ENGINE_PARAMS} ${APPEND_ARGS}"
        MODIFIED_PARAMS=true
      fi
      if [[ -n "$PREPEND_ARGS" ]] ; then
        ENGINE_PARAMS="${PREPEND_ARGS} ${ENGINE_PARAMS}"
        MODIFIED_PARAMS=true
      fi

      if [[ "$MODIFIED_PARAMS" = true ]] ; then
        echo "Updated engine params are ${ENGINE_PARAMS}"
      fi
}

function cleanup() {
    set +e
    wait $ENGINE_PID
    set -e
    echo "=== Running Cleanup ==="

    echo "=== Bye! ==="
}

function cancel() {
    # AWS Batch sends a SIGTERM to a container if its job is cancelled/terminated
    # forward this signal to engine so that it can cancel any pending workflow jobs

    set +e  # ignore errors here
    echo "=== !! CANCELLING WORKFLOW !! ==="
    echo "stopping ${ENGINE_NAME} pid: $ENGINE_PID"
    kill -TERM "$ENGINE_PID"
    echo "waiting .."
    wait $ENGINE_PID
    echo "=== !! cancellation complete !! ==="
    set -e
}

trap "cancel; cleanup" TERM
trap "cleanup" EXIT


# AWS Batch places multiple jobs on an instance
# To avoid file path clobbering use the JobID
# to create a unique path. This is important if /opt/work
# is mapped to a filesystem external to the container
GUID="$AWS_BATCH_JOB_ID"

if [ "$GUID" = "/" ]; then
    GUID=`date | md5sum | cut -d " " -f 1`
fi

# Make the directory we will work in
mkdir -p /mnt/efs/snakemake/$GUID
cd /mnt/efs/snakemake/$GUID

if [[ "$ENGINE_PROJECT" =~ ^s3://.* ]]; then
    echo "== Staging S3 Project =="
    ENGINE_PROJECT_DIRECTORY="."
    echo "Copying from ${ENGINE_PROJECT} to '${ENGINE_PROJECT_DIRECTORY}/'"
    aws s3 cp $ENGINE_PROJECT "${ENGINE_PROJECT_DIRECTORY}/"
    find $ENGINE_PROJECT_DIRECTORY -name '*.zip' -execdir unzip -o '{}' ';'
    ls -l $ENGINE_PROJECT_DIRECTORY
    MANIFEST_JSON=${ENGINE_PROJECT_DIRECTORY}/MANIFEST.json
    if test -f "$MANIFEST_JSON"; then
     handleManifest
    else
      ENGINE_PROJECT="${ENGINE_PROJECT_DIRECTORY}"
      ENGINE_PARAMS="${ENGINE_PARAMS} --aws-batch-tags AWS_BATCH_PARENT_JOB_ID=${AWS_BATCH_JOB_ID}  --aws-batch-efs-project-path=snakemake/$GUID --latency-wait 60"
    fi
fi
echo "== Finding the project in  =="
echo "cd ${ENGINE_PROJECT}"
cd ${ENGINE_PROJECT}
echo "== Running Workflow =="
echo "${ENGINE_RUN_CMD} ${ENGINE_PARAMS}"
$ENGINE_RUN_CMD $ENGINE_PARAMS & ENGINE_PID=$!

echo "${ENGINE_NAME} pid: $ENGINE_PID"
jobs
echo "waiting .."
wait $ENGINE_PID
