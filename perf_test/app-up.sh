#!/bin/bash
set -x
#sla
date +%Y%m%d%H%M%S

kubectl  create -f nginx.yaml
#kubectl  create -f alpine.yaml 

kubectl  create -f nginx1.yaml
#kubectl  create -f alpine1.yaml 

kubectl  create -f nginx2.yaml
#kubectl  create -f alpine2.yaml 

kubectl  create -f nginx3.yaml
#kubectl  create -f alpine3.yaml 

kubectl  create -f nginx4.yaml
#kubectl  create -f alpine4.yaml 

#non-sla
#kubectl  create -f alpine.yaml

