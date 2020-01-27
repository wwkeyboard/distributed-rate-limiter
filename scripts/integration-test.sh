#!/bin/bash
set -e

# Bounce the redis server if it's not running.  This gets us a known starting environment.  We also do this because if the server is running it'll cause start-redis-server.sh to throw an error, stoping this script from running. This will happen if the build doesn't return 0, which will stop the script from nuking the Redis server and destroying evidence of what when wrong.

./scripts/stop-redis-server.sh

./scripts/start-redis-server.sh

set -x

go test -tags=integration ./...

set +x

./scripts/stop-redis-server.sh
