#!/usr/bin/env bash

HOME=`pwd`
HOME_DATA=$PWD"/data/ssdb"

for port in 7027 7028 7029; do
    data_dir=$HOME_DATA/ssdb_$port
    conf_file=$data_dir/ssdb.conf

    ssdb-server $conf_file -s stop
done
