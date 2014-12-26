#!/bin/bash

./thingz-agent -source=demo \
               -strategy=cpu:3,cpus:5,mem:4,swap:5,load:5 \
               -publisher="http://thingz:thingz@localhost:8086/thingz" \
               -verbose=true
