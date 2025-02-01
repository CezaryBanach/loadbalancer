#!/bin/bash

echo "starting 3 servers"
./bin/backend 8081 & echo $! >> backend_pids.txt
./bin/backend 8082 & echo $! >> backend_pids.txt
./bin/backend 8083 & echo $! >> backend_pids.txt

echo "All servers started"
