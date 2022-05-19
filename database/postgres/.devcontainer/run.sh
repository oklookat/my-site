#!/bin/bash

initdb -D $PGDATA

# allow all ip/ipv6 connections
echo "host    all             all             0.0.0.0/0            md5" >> $PGDATA/pg_hba.conf
echo "host    all             all             ::0/0                md5" >> $PGDATA/pg_hba.conf

# init database & user
pg_ctl start -D $PGDATA
psql -c "CREATE DATABASE elven"
psql -c "ALTER USER postgres WITH PASSWORD 'root';"
pg_ctl stop -D $PGDATA

# run
postgres -c config_file=/app/postgres.conf