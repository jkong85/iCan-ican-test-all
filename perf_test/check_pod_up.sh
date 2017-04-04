#!/bin/bash

until [ $(kubectl get pods | grep -o 'Running' | wc -l) -eq $1 ]
do
    echo "wait for one seconds"
    sleep 0.1
done

