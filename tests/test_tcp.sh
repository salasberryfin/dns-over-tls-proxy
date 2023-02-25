#!/bin/bash

while :
do
    dig @127.0.0.1 -p 5353 google.com +tcp
    sleep 1
done

