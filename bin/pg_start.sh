#!/bin/bash

HOME=/mnt

DATADIR=$HOME/data/pgsql

pg_ctl start -D $DATADIR -l $DATADIR/postgresql.log
