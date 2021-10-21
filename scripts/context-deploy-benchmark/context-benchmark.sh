#!/usr/bin/env bash

CONTEXT_AMOUNT=1

run_contexts() {
  for ((i = 1; i <= CONTEXT_AMOUNT; i++))
  do
    agc context deploy -c "$1$i";
  done
}

get_average_runtime() {
  START_SECONDS=$(date +%s)
  run_contexts $1
  END_SECONDS=$(date +%s)
  DIFF_SECONDS=$((END_SECONDS-START_SECONDS))
  AVERAGE_SECONDS=$(((DIFF_SECONDS%60)/CONTEXT_AMOUNT))
  AVERAGE_MINUTES=$(((DIFF_SECONDS/60)/CONTEXT_AMOUNT))
  echo "Average $1 deployment time across $CONTEXT_AMOUNT contexts is: $AVERAGE_MINUTES m $AVERAGE_SECONDS s"
}

validate_and_run() {
  if [ $1 != "1" ] && [[ $2 != "cromwell" || $2 != "nextflow" ]]
  then
    echo "Please specify engine type to benchmark (cromwell | nextflow)"
  else
    get_average_runtime $2
  fi
}

validate_and_run $# $1
