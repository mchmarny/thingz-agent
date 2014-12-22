#!/bin/bash

./thingz -source=laptop \
         -strategy=cpu:1,mem:1,swap:5,load:5 \
         -publisher="http://test:test@127.0.0.1:8086/thingz"