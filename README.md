# thingz-agent

> Thingz understood

This agent works in tandem with the [thingz-server](https://github.com/mchmarny/thingz-server) to provide demonstration of both the dynamic modeling to support actuation as well as forensic query and visualization.

## Install

Once you have golang [installed](http://golang.org/doc/install):

```
go get github.com/mchmarny/thingz-agent
go install github.com/mchmarny/thingz-agent
```


## Run

Export a few environment variables

* `THINGZ_HOST` is the [thingz-server](https://github.com/mchmarny/thingz-server) host
* `THINGZ_SECRET` is the agent secret 

Once these two variables are defined, you can imply start the agent using the following command: 


```
./thingz-agent --source="${HOSTNAME}" \
               --strategy=cpu:10,cpus:60,mem:15,swap:30,load:15 \
               --publisher="http://agent:${THINGZ_SECRET}@${THINGZ_HOST}:8086/thingz" \
               --verbose=true
```                  

> Note, the `source` parameter can be anything that uniquely identifies this host. If the `HOSTNAME` includes periods (`.`) you can override this value with your own nickname for this host (e.g. `server-1`)

## Strategies 

Strategies tell the agent what dimensions and how often it should report to the `thingz-server`. You specify a strategy by defining one or more of Key/Value pares where the key is the dimension and the value is the reporting frequency in seconds.

> Available Dimensions:


* **cpu**  - user, nice, sys, idle, wait, total
* **cpus** - just like above but for each CPU individually
* **mem**  - free, used, actual free, actual used, total
* **swap** - free, used, total
* **load** - avg for 1, 5, 15 min
* **proc** - process name and resident memory

