#!/bin/bash

set -ue

IP=$(ifconfig en0 | grep inet | grep broadcast | cut -d ' ' -f 2)

echo 本机IP: $IP

# TOKEN=docker-etcd
# CLUSTER_STATE=new
# NAME_1=node-1
# NAME_2=node-2
# NAME_3=node-3
HOST_1=$IP
HOST_2=$IP
HOST_3=$IP

# CLUSTER=${NAME_1}=http://${HOST_1}:1001,${NAME_2}=http://${HOST_2}:2001,${NAME_3}=http://${HOST_3}:3001

ENDPOINTS=$HOST_1:1000,$HOST_2:2000,$HOST_3:3000

# etcdctl --endpoints="$ENDPOINTS" member list

# etcdctl --endpoints="$ENDPOINTS" put project "leekbox"
# etcdctl --endpoints="$ENDPOINTS" put token "leooo"

# etcdctl --endpoints="$ENDPOINTS" --write-out="json" get project
etcdctl --endpoints="$ENDPOINTS" get q1mi

etcdctl --write-out=table --endpoints="$ENDPOINTS" endpoint status
