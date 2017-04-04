#!/bin/bash

ssh ican-1 sudo tc qdisc show dev br
ssh ican-1 sudo tc class show dev br
ssh ican-1 sudo tc filter show dev br
