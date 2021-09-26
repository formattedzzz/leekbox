#!/bin/bash
set -ue

IP=$(ifconfig en0 | grep inet | grep broadcast | cut -d ' ' -f 2)

if [ "$IP" ]; then
  echo "本机IP:$IP"
else
  IP=$(ifconfig en4 | grep inet | grep broadcast | cut -d ' ' -f 2)
  echo "本机IP:$IP"
fi

sed -e "s/localhost/$IP/g" docker-compose-base.yml >docker-compose.yml

docker-compose up -d
