#!/bin/bash

DIR="$(pwd)"

$DIR/thingz-agent --source=demo \
                  --strategy=cpu:3,cpus:5,mem:4,swap:5,load:5 \
                  --publisher="stdout" \
                  --verbose=true
