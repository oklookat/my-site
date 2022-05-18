#!/bin/bash

initdb -D $PGDATA
pg_ctl start -D $PGDATA
psql -c "CREATE DATABASE elven"
psql -c "ALTER USER postgres WITH PASSWORD 'root';"