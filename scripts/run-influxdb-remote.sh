#!/bin/bash

DIR="$(pwd)"

$DIR/thingz-agent --source="${HOSTNAME}" \
                  --strategy=cpu:10,cpus:60,mem:15,swap:30,load:15 \
                  --publisher="influxdb" \
                  --publisher-args="http://agent:${THINGZ_SECRET}@${THINGZ_HOST}:8086/thingz" \
                  --verbose=true