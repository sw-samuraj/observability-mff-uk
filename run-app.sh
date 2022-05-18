#!/bin/sh -eu

echo "Starting observability app..."
echo "Logs are stored in the _logs/observability.log file."
echo "You can kill the app with Ctrl-C."
./observability -d "http://localhost:5050"
