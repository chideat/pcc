#!/bin/bash

HOME=/db

DATADIR=$HOME/data/pgsql

pg_ctl stop -D $DATADIR
