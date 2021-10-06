#!/usr/bin/env bash
set -e

export PGPASSWORD=test
# TODO добавить создание схем
psql -U program -d services <<-EOSQL

EOSQL