#!/bin/bash

HOME=/mnt

DATADIR=$HOME/data/pgsql
mkdir -p $DATADIR

pg_ctl stop -D $DATADIR
