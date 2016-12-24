#!/usr/bin/env bash

HOME=`pwd`
HOME_DATA=$PWD"/data/redis"

mkdir -p $HOME_DATA/redis_conf
mkdir -p $HOME_DATA/redis_pids
mkdir -p $HOME_DATA/redis_logs

port=7001

data_dir=$HOME_DATA/redis_$port
mkdir -p $data_dir
pid_path=$HOME_DATA/redis_pids/redis_$port.pid
log_path=$HOME_DATA/redis_logs/redis_$port.log
conf_path=$HOME_DATA/redis_conf/redis_$port.conf

sed -e "s|__PORT__|$port|" \
    -e "s|__DATA_DIR__|$data_dir|" \
    -e "s|__PID_PATH__|$pid_path|" \
    -e "s|__LOG_PATH__|$log_path|" \
    $HOME/conf/redis.conf > $conf_path

echo "start redis at:" [$data_dir] "port:" $port
echo "redis-server $conf_path"
redis-server $conf_path
