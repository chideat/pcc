#!/usr/bin/env bash

HOME="/ssdb/data"

for port in 6379; do
    data_dir=$HOME/ssdb_$port
    conf_file=$data_dir/ssdb.conf

    ssdb-server $conf_file -s stop
done
