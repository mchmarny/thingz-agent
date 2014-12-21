#!/bin/bash

go build

./thingz -verbose=true \
         -source=laptop \
         -strategy=cpu:5,mem:10,swap:60,load:5