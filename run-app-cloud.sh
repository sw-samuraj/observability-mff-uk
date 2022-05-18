#!/bin/sh -eu

echo "Starting observability app..."
echo "Logs are stored in the _logs/observability.log file."
echo "You can kill the app with Ctrl-C."
./observability -d "http://service1.edu.dobias.info:5050" -t "http://grafana.edu.dobias.info:14268/api/traces"
# ./observability -p 5050 -n downstream-1 -d "http://service2.edu.dobias.info:6060" -t "http://grafana.edu.dobias.info:14268/api/traces"
# ./observability -p 6060 -n downstream-2 -t "http://grafana.edu.dobias.info:14268/api/traces"
