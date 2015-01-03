#!/bin/bash

# update
sudo apt-get update -y

# install git
sudo apt-get install git mercurial -y

# Go
if [ -d go ]; then
  echo "Go already installed"
else
  wget -q https://storage.googleapis.com/golang/go1.4.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.4.linux-amd64.tar.gz
fi

# GOPATH
echo '# go
export GOROOT=/usr/local/go
export GOPATH=/home/ubuntu/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' | sudo tee -a /etc/profile

# Finally source environment variables from the updated profile
source /etc/profile

# Dependencies
go get github.com/cloudfoundry/gosigar
go get github.com/influxdb/influxdb
go get github.com/Shopify/sarama

# mchmarny
mkdir ~/go/src/github.com/mchmarny


