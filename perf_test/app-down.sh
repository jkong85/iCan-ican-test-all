#!/bin/bash
set -x
#sla
kubectl  delete -f nginx.yaml
kubectl  delete -f alpine.yaml 

kubectl  delete -f nginx1.yaml
kubectl  delete -f alpine1.yaml 

kubectl  delete -f nginx2.yaml
kubectl  delete -f alpine2.yaml 

kubectl  delete -f nginx3.yaml
kubectl  delete -f alpine3.yaml 

kubectl  delete -f nginx4.yaml
kubectl  delete -f alpine4.yaml 
#non-sla
#kubectl  create -f alpine.yaml

