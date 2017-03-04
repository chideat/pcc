#!/bin/bash

HOME=/mnt

DATADIR=$HOME/data/pgsql
mkdir -p $DATADIR

pg_ctl start -D $DATADIR -l $DATADIR/postgresql.log
