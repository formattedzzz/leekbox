#!/bin/bash

set -eu

source ./init-db.sh

echo main-sh

echo $HOME

leo=" "

if [ -n "$leo" ]; then
  echo "len not 0"
else
  echo "len is 0"
fi

if [ $leo ]; then
  echo "not null"
else
  echo null
fi

LEO=liufulin && echo $LEO