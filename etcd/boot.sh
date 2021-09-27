#!/bin/bash
set -ue

IP=$(ifconfig en0 | grep inet | grep broadcast | cut -d ' ' -f 2)

if [ "$IP" = "" ]; then
  echo "en0不存在" # 取ecs上面的
  # grep 正则很多系统上参数不同 有的是-E 有的是-P 转义的兼容性好一点
  IP=$(ifconfig eth0 | grep inet | grep -o '\([0-9]\+.\)\+' | sed -n '1p')
fi

echo "本机IP:$IP"

TOKEN=leekbox

sed -e "s/\${LOCIP}/$IP/g" -e "s/\${TOKEN}/$TOKEN/g" docker-compose-base.yml >docker-compose.yml

docker-compose up -d
