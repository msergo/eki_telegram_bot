#!/bin/sh
set -e

if [ "${CACHE_ENGINE}" != "sqlite" ]; then
  echo "Waiting for Redis to be available"
  while ! nc -z redis 6379; do sleep 2; done
  echo "Redis connected. Starting app"
else
  echo "Cache engine: SQLite. Skipping Redis check. Starting app"
fi

$1
exit $?
