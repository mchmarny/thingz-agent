#!/bin/bash

./thingz-agent -source=demo \
               -strategy=cpu:3,cpus:5,mem:3,swap:10,load:15 \
               -publisher="http://thingz:thingz@localhost:8086/thingz"