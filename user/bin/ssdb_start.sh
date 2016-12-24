#!/usr/bin/env bash

HOME=`pwd`
HOME_DATA=$PWD"/data/ssdb"

for port in 7027 7028 7029; do
    data_dir=$HOME_DATA/ssdb_$port
    mkdir -p $data_dir

    pid_file=$data_dir/ssdb.pid
    log_file=$data_dir/log.txt
    conf_file=$data_dir/ssdb.conf
    
    sed -e "s|__PORT__|$port|" \
        -e "s|__DATA_DIR__|$data_dir|" \
        -e "s|__PID_FILE__|$pid_file|" \
        -e "s|__LOG_FILE__|$log_file|" \
        $HOME/conf/ssdb.conf> $conf_file
    
    ssdb-server -d $conf_file
done
