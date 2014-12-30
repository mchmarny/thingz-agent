#!/bin/bash

DIR="$(pwd)"

$DIR/thingz-agent --source="${HOSTNAME}" \
                  --strategy=cpu:10,cpus:60,mem:15,swap:30,load:15 \
                  --publisher="kafka" \
                  --publisher-args="thingz,localhost:9092" \
                  --verbose=true