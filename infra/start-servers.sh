#!/usr/bin/env bash

LOG_DIR=/var/log/axolotl/

MLFLOW_PORT=8090
MLFLOW_LOG=$LOG_DIR/mlflow.log

JUPYTER_PORT=8080
JUPYTER_LOG=$LOG_DIR/jupyter.log

AIRFLOW_PORT=8070
AIRFLOW_WEB_LOG=$LOG_DIR/airflow-api-server.log
AIRFLOW_SCHEDULER_LOG=$LOG_DIR/airflow-scheduler.log

nohup mlflow server --host 0.0.0.0 --port $MLFLOW_PORT > $MLFLOW_LOG 2>&1 &

nohup jupyter lab --allow-root --ip 0.0.0.0 --port $JUPYTER_PORT > $JUPYTER_LOG 2>&1 &

nohup airflow api-server -p $AIRFLOW_PORT > $AIRFLOW_WEB_LOG 2>&1 &

nohup airflow scheduler > $AIRFLOW_SCHESDULER_LOG 2>&1 &

tail -f $MLFLOW_LOG $JUPYTER_LOG $AIRFLOW_WEB_LOG $AIRFLOW_SCHEULER_LOG
