#!/bin/bash

DIR="$(pwd)"

$DIR/thingz-agent --source=demo \
                  --strategy=cpu:1,cpus:1,mem:1,swap:5,load:5 \
                  --publisher="influxdb" \
                  --publisher-args="udp://thingz:thingz@localhost:4444/thingz" \
                  --verbose=true
