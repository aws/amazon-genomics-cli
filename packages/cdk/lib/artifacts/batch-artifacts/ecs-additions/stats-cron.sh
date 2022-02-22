#!/bin/bash

echo "Running docker stats script"

REGION=$(curl http://169.254.169.254/latest/meta-data/placement/region 2> /dev/null)
INSTANCE_ID=$(curl http://169.254.169.254/latest/meta-data/instance-id 2> /dev/null)
STATS=$(docker stats --no-stream --format "{\"container\": \"{{ .Container }}\",\"name\": \"{{ .Name }}\", \"memory\": { \"raw\": \"{{ .MemUsage }}\", \"percent\": \"{{ .MemPerc }}\"}, \"cpu\": \"{{ .CPUPerc }}\"}" | jq '.' -s -c )
NUM_CONTAINERS=$(echo "$STATS" | jq '. | length')

for (( i=0; i<$NUM_CONTAINERS; i++ )) 
do CPU=$(echo "$STATS" | jq -r .[$i].cpu | sed 's/%//')
MEMORY=$(echo "$STATS" | jq -r .[$i].memory.percent | sed 's/%//')
RAW_MEM=$(echo "$STATS" | jq -r .[$i].memory.raw | sed 's/%//')
CONTAINER=$(echo $STATS | jq -r .[$i].container)
CONTAINER_NAME=$(echo $STATS | jq -r .[$i].name)
echo "emitting stats to cw"
aws cloudwatch put-metric-data --metric-name CPU --namespace DockerStats --unit Percent --value $CPU --dimensions InstanceId=$INSTANCE_ID,ContainerId=$CONTAINER,ContainerName=$CONTAINER_NAME --region $REGION
aws cloudwatch put-metric-data --metric-name Memory --namespace DockerStats --unit Percent --value $MEMORY --dimensions InstanceId=$INSTANCE_ID,ContainerId=$CONTAINER,ContainerName=$CONTAINER_NAME --region $REGION
aws cloudwatch put-metric-data --metric-name RawMemory --namespace DockerStats --unit Count --value $RAW_MEM --dimensions InstanceId=$INSTANCE_ID,ContainerId=$CONTAINER,ContainerName=$CONTAINER_NAME --region $REGION

done