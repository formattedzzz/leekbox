#!/bin/bash
set -e

current_dir=$(pwd)
echo "当前目录:$current_dir"

WORK_DIR=${1:-$GOPATH}

echo "工作目录(该项目所在的目录):$WORK_DIR"

PB_DIR="$WORK_DIR/leekbox/pb"

protoc -I="$PB_DIR" --go_out="$WORK_DIR" "$PB_DIR"/*.proto

echo "编译完成！"