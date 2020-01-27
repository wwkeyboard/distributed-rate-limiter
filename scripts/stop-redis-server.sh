#!/bin/bash

set -x
# This will probably remove the container, if it was started with start-redis-server.sh
docker stop some-redis
set +x
exit 0