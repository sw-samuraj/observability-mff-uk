#!/bin/bash -eu

LOG_DIR="_logs"
LOG_FILE="${LOG_DIR}/observability.log"
POSITION_FILE="${LOG_DIR}/positions.yaml"

kill_app () {
  APP_PID=$(pgrep "${APP}" || echo "" )
  if [ -z "${APP_PID}" ]
  then
    echo "${APP} is not running."
  else
    echo "Killing ${APP} running with PID ${APP_PID}."
    kill -9 "${APP_PID}"
  fi
}

kill_downstream () {
  DOWNSTREAM=$(pgrep -a observability | grep downstream || echo "")
  if [ -z "${DOWNSTREAM}" ]
  then
    echo "Downstream services are not running."
  else
    echo "Killing downstream services."
    pkill -f 'observability.*downstream'
  fi
}

# TODO Modify following process names according your environment.
# Following app names are symlinks - if you don't have them you should either
# create them or change variables to real process names in your environment.
APP="loki"
kill_app
APP="promtail"
kill_app
APP="prometheus"
kill_app
APP="jaeger"
kill_app
APP="grafana"
kill_app

kill_downstream

if [ -f "${LOG_FILE}" ]
then
  echo "Deleting old log files..."
  find "${LOG_DIR}" -name '*.log' -delete
fi

if [ -f "${POSITION_FILE}" ]
then
  echo "Deleting old position file..."
  rm "${POSITION_FILE}"
fi
