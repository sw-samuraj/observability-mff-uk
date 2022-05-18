#!/bin/bash -eu

MAX_REQ=12
URL="localhost:4040/"
# URL="student.edu.dobias.info:4040/"
TIMESTAMP=$(date +%s)

COUNTER=1
while [ "$COUNTER" -le $MAX_REQ ]
do
  echo "------------ Load request #${COUNTER} ------------"
  echo ""
  curl -v -H "X-Request-ID: load-${COUNTER}-${TIMESTAMP}" "${URL}"
  echo ""
  echo ""
  ((COUNTER++))
done
