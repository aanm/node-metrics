#!/usr/bin/env sh
cat /proc/meminfo | grep "$1" | awk '{ print $2 }'