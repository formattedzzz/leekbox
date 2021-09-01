#!/bin/bash

set -eu

echo "首次启动mysql镜像 执行初始化脚本 数据库列表:"

mysql -uroot -p$MYSQL_ROOT_PASSWORD <<EOF
show databases;
EOF

echo "数据卷目录：$LEEKBOX_HOME/mysql-lib"
