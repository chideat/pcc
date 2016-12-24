#!/bin/bash

for i in "7021" "7022" "7023" "7024" "7025" "7026"; do
    port=$i
    echo "stop redis" [$i]": redis-cli -p $port shutdown"
    redis-cli -p $port shutdown 
done
