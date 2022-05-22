#!/bin/bash

# init
initdb -D ${PGDATA}

# allow all ip/ipv6 connections
echo "host all all 0.0.0.0/0 md5" >> ${PGDATA}/pg_hba.conf
echo "host all all ::0/0 md5" >> ${PGDATA}/pg_hba.conf

# create DB & change user password
pg_ctl start -D ${PGDATA}
psql -c "CREATE DATABASE ${POSTGRES_DB}"
psql -c "ALTER USER ${POSTGRES_USER} WITH PASSWORD '${POSTGRES_PASSWORD}';"
pg_ctl stop -D ${PGDATA} -m fast

exec "$@"