#!/bin/bash

go build

./thingz -verbose=true \
         -source=laptop \
         -strategy=cpu:3,mem:3,swap:60,load:5