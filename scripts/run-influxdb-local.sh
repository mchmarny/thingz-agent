#!/bin/bash

DIR="$(pwd)"

$DIR/thingz-agent --source=mbp13 \
                  --strategy=cpu:1,cpus:1,mem:1,swap:5,load:5 \
                  --publisher="influxdb" \
                  --publisher-args="http://thingz:thingz@localhost:8086/thingz" \
                  --verbose=true
