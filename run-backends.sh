#!/bin/sh -eu

LOKI_HOME="${HOME}/dev/loki"
PWD=$(pwd)

echo "Starting Loki in the background..."
"${LOKI_HOME}/loki-linux-amd64" -config.file "${PWD}/_config/loki/loki-local-config.yaml" &

echo "Starting Promtail in the background..."
"${LOKI_HOME}/promtail-linux-amd64" -config.file "${PWD}/_config/loki/promtail-local-config.yaml" &
