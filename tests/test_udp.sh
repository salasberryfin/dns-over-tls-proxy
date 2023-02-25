#!/bin/bash

while :
do
    dig @127.0.0.1 -p 5354 google.com
    sleep 1
done
