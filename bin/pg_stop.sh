#!/bin/bash

HOME=/mnt

DATADIR=$HOME/data/pgsql

pg_ctl stop -D $DATADIR
