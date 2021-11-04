#!/usr/bin/env bash

CONTEXT_AMOUNT=5
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

deploy_contexts_sequential() {
  for ((i = 1; i <= CONTEXT_AMOUNT; i++))
  do
    agc context deploy -c "$1$i";
  done
}

cleanup_contexts() {
    for ((i = 1; i <= CONTEXT_AMOUNT; i++))
    do
      agc context destroy -c "$1$i";
    done
}

deploy_contexts_all() {
  agc context deploy --all;
}

get_average_runtime() {
  START_SECONDS=$(date +%s)
  $1 $2
  END_SECONDS=$(date +%s)
  DIFF_SECONDS=$((END_SECONDS-START_SECONDS))
  AVERAGE_SECONDS=$(((DIFF_SECONDS%60)/CONTEXT_AMOUNT))
  AVERAGE_MINUTES=$(((DIFF_SECONDS/60)/CONTEXT_AMOUNT))
  echo "Average $2 $3 time across $CONTEXT_AMOUNT contexts is: $AVERAGE_MINUTES m $AVERAGE_SECONDS s"
}

get_all_averages() {
      get_average_runtime deploy_contexts_sequential "$1" "sequential deployment"
      cleanup_contexts "$1"
      get_average_runtime deploy_contexts_all "$1" "all deployment"
      get_average_runtime cleanup_contexts "$1" "destruction"
}

validate_and_run() {
  if [ $1 != "1" ] || [[ $2 != "cromwell" && $2 != "nextflow" ]]
  then
    echo "Please specify engine type to benchmark (cromwell | nextflow)"
  else
    cd "$SCRIPT_DIR/$2"
    get_all_averages $2
  fi
}

validate_and_run $# $1
