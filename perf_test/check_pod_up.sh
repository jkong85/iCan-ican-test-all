#!/bin/bash

echo "wait for all pods up"
until [ $(kubectl get pods | grep -o 'Running' | wc -l) -eq $1 ]
do
    sleep 2 
done

