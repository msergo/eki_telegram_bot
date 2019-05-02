#!/bin/sh
set -e
echo Starting redis
echo Command to run: $1
redis-server & $1
exit $1