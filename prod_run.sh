# !/bin/bash
set -m
./main -region=$REGION_ENV -secret=$SECRET -deployment=$DEPLOYMENT &
./filebeat/filebeat -e

fg %1