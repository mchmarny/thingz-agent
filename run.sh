#!/bin/bash

go build

./thingz -verbose=true \
         -source=laptop \
         -strategy=cpu:1