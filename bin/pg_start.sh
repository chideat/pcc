#!/bin/bash

HOME=/db

DATADIR=$HOME/data/pgsql

pg_ctl start -D $DATADIR -l $DATADIR/postgresql.log
