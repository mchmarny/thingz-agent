#!/bin/bash

DIR="$(pwd)"

$DIR/thingz-agent --source="${HOSTNAME}" \
                  --strategy=cpu:3,cpus:60,mem:3,swap:30,load:5 \
                  --publisher="influxdb" \
                  --publisher-args="udp://agent:${THINGZ_SECRET}@${THINGZ_HOST}:4444/thingz" \
                  --verbose=true