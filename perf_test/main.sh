#!/bin/bash

date +%Y%m%d%H%M%S
./app-up.sh
./check_pod_up.sh 195
date +%Y%m%d%H%M%S
#./tc 0 1 s
./tc 5 1 s
date +%Y%m%d%H%M%S

#./app-down.sh
#./tc 0 50 d
