#!/usr/bin/env bash
set -e

# TODO для корректного создания схем в Postgres прописать свой вариант
export VARIANT=v1
export SCRIPT_PATH=/docker-entrypoint-initdb.d/
export PGPASSWORD=test
psql -U program -d services -f "$SCRIPT_PATH/schemes/schema-$VARIANT.sql"
