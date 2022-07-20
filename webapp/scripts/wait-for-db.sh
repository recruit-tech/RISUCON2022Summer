#!/bin/bash
# wait-for-it.sh

set -e

until nc -z $MYSQL_HOST $MYSQL_PORT; do
  >&2 echo "mysql is unavailable - sleeping"
  sleep 3
done
>&2 echo "mysql is up - executing command"

exec $@
