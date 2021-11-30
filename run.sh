#!/bin/bash

set -e

echo "ENV"
env

echo "LS1"
ls -la

echo "LS2"
ls -la /usr/local/bin

echo "1"
echo litestream restore -o $DB_PATH $REPLICA_URL
litestream restore -if-replica-exists -o $DB_PATH $REPLICA_URL

echo "2"
echo litestream replicate --exec app $DB_PATH $REPLICA_URL
litestream replicate --exec app $DB_PATH $REPLICA_URL
