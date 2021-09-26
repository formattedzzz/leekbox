#!/bin/bash

# IP=$(ifconfig en0 | grep inet | grep broadcast | cut -d ' ' -f 2)

# CLUSTER=node1=http://"$IP"

# CLUSTER_TOKEN=leekbox

# echo 本机IP："$IP"

# docker run \
#   -p 2379:2379 \
#   -p 2380:2380 \
#   --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
#   --name etcd \
#   quay.io/coreos/etcd:v3.5.0 \
#   etcd \
#   --name node1 \
#   --data-dir /etcd-data \
#   --listen-client-urls http://"$IP":2379 \
#   --advertise-client-urls http://"$IP":2379 \
#   --listen-peer-urls http://"$IP":2380 \
#   --initial-cluster "$CLUSTER" \
#   --initial-cluster-token "$CLUSTER_TOKEN" \
#   --initial-cluster-state new \
#   --log-level info \
#   --log-outputs stderr

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
