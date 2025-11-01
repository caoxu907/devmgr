#!/bin/bash

EMQTT_BENCH="/data/pc/emq-test/bin/emqtt_bench"
TOPIC="cloud/device_key_16/telemetry"
HOST="emqx-node1.cluster.local"
USER="test"
PASS="password@123"
WORKERS=1
INTERVAL=10000   # 毫秒

SUB_DEV="subdevice_key_16"
POINT1="point_key_16_20"
POINT2="point_key_16_18"

while true; do
    TIME=$(date +%s%3N)
    VALUE1=$((RANDOM % 20 + 80))
    VALUE2=$(awk -v r=$RANDOM 'BEGIN{srand(r); printf "%.1f", 1+rand()*99}')
    JSON="{\"payloads\":[{\"deviceid\":\"$SUB_DEV\",\"pointid\":\"$POINT1\",\"value\":$VALUE1,\"time\":$TIME},{\"deviceid\":\"$SUB_DEV\",\"pointid\":\"$POINT2\",\"value\":$VALUE2,\"time\":$TIME}]}"
    timeout -s KILL 5 $EMQTT_BENCH pub -t $TOPIC -h $HOST -c $WORKERS -I $INTERVAL -u "$USER" -P "$PASS" -m "$JSON"
    sleep 10
done

