#!/bin/sh
set -e
echo Waiting for Redis to be available
while ! nc -z redis 6379; do sleep 2; done
echo "Redis connected. Starting app"
$1
exit $1