#!/bin/bash

date +%Y%m%d%H%M%S
#../app-up.sh
./check_pod_up.sh $1
date +%Y%m%d%H%M%S
./tc $2 1 s
date +%Y%m%d%H%M%S

