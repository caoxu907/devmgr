#!/bin/bash

# 运行参数o
CLIENT_ID="device_key_2"
LOG_FILE="edge_test.log"
BIN="./test_mqtt_edge.test"

# 检查可执行文件是否存在
if [ ! -f "$BIN" ]; then
  echo "错误: 未找到 $BIN"
  exit 1
fi

# 后台运行并重定向日志
nohup $BIN -test.v -id="$CLIENT_ID" > "$LOG_FILE" 2>&1 &

# 输出进程信息
PID=$!
echo "$BIN 已在后台运行，PID: $PID"
echo "日志输出: $LOG_FILE"
