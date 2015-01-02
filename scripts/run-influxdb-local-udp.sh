#!/bin/bash

DIR="$(pwd)"

cd $DIR

./thingz-agent --source=mbp13 \
               --strategy=cpu:1,cpus:1,mem:1,swap:5,load:5 \
               --publisher="influxdb" \
               --publisher-args="udp://thingz:thingz@localhost:4444/thingz" \
               --verbose=true
