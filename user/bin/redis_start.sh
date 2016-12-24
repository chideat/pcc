#!/usr/bin/env bash

HOME=`pwd`
HOME_DATA=$PWD"/data/redis"

mkdir -p $HOME_DATA/redis_conf
mkdir -p $HOME_DATA/redis_pids
mkdir -p $HOME_DATA/redis_logs

for i in "7021" "7022" "7023" "7024" "7025" "7026"; do
    data_dir=$HOME_DATA/redis_$i
    mkdir -p $data_dir
    pid_path=$HOME_DATA/redis_pids/redis_$i.pid
    log_path=$HOME_DATA/redis_logs/redis_$i.log
    conf_path=$HOME_DATA/redis_conf/redis_$i.conf

    port=$i
    sed -e "s|__PORT__|$port|" \
        -e "s|__DATA_DIR__|$data_dir|" \
        -e "s|__PID_PATH__|$pid_path|" \
        -e "s|__LOG_PATH__|$log_path|" \
        $HOME/conf/redis.conf > $conf_path

    echo "start redis" [$i] "at:" [$data_dir] "port:" $port
    echo "redis-server $conf_path"
    redis-server $conf_path
done
