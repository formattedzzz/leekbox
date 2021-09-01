#!/bin/bash

set -e

mkdir -p $HOME/.leekbox

LEEKBOX_HOME=$HOME/.leekbox

init=false

for param in "$@"; do
  if [ $param == '--init' ]; then
    init=true
  fi
done

if [ "$init" == true ]; then
  echo '格式化除宿主机数据库数据卷'
  volumns=("mysql-etc" "mysql-lib" "mysql-files")
  for file in ${volumns[*]}; do
    echo "$file"
    rm -rf $LEEKBOX_HOME/$file
  done
fi

docker run \
  -e MYSQL_ROOT_PASSWORD=leekbox \
  -e LEEKBOX_HOME=$LEEKBOX_HOME \
  -v $LEEKBOX_HOME/mysql-lib:/var/lib/mysql \
  -v $LEEKBOX_HOME/mysql-etc:/etc/mysql \
  -v $LEEKBOX_HOME/mysql-files:/var/lib/mysql-files \
  -p 33061:3306 \
  mysql:leekbox \
  --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
