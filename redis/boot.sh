#!/bin/bash

# docker run --name redis6 -d -p 6379:6379 -v $(pwd):/data redis:6.0 redis-server --appendonly yes

docker volume create --driver local \
  --opt type=nfs \
  --opt o=addr=127.0.0.1,rw \
  --opt device=:/Users/liufulin/Desktop/volumes \
  leooo

docker volume create --driver local \
  --opt type=btrfs \
  --opt device=/dev/sda2 \
  foo
