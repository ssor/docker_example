#!/usr/bin/env bash

docker run \
        -v /Users/zhangquanzhi/docker_vol:/config \
        --name host$1 \
        --network net_on_mac \
        -it -d -p  $1:$1 hello:3 \
         --port=$1
