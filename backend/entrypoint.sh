#!/bin/sh
# ./backend/entrypoint.sh

set -e

echo "Backend Entrypoint: Starting data synchronization task..."

/app/stockify_datasync

echo "Backend Entrypoint: Data synchronization task finished."
echo "Backend Entrypoint: Starting main backend application..."

exec "$@"
