#!/bin/bash

for i in "7001"; do
    port=$i
    echo "stop redis" [$i]": redis-cli -p $port shutdown"
    redis-cli -p $port shutdown 
done

