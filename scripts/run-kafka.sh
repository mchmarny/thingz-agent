#!/bin/bash

DIR="$(pwd)"

$DIR/thingz-agent --source="${HOSTNAME}" \
                  --strategy=cpu:1,cpus:5,mem:1,swap:3,load:5 \
                  --publisher="kafka" \
                  --publisher-args="thingz,localhost:9092" \
                  --verbose=true