#!/bin/bash

HOSTPARAMS="--host cockroachdb --insecure"
SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"

$SQL -e "CREATE USER IF NOT EXISTS netra;"
$SQL -d netra -e "CREATE DATABASE IF NOT EXISTS netra WITH ENCODING = 'UTF8';"
$SQL -d netra -e "GRANT ALL ON DATABASE netra TO netra;"
$SQL -d netra -e "CREATE TABLE IF NOT EXISTS issues (id SERIAL PRIMARY KEY, title VARCHAR NOT NULL, description VARCHAR, priority int);"
