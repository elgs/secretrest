#!/bin/bash

go get -u github.com/elgs/secretrest
cd /root/secretrest
./shutdown.sh
nohup secretrest &
