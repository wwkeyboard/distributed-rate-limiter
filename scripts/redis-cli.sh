#!/bin/sh

# Opens a console to the development Redis server started with ./redis-server.sh

docker run -it --network some-network --rm redis redis-cli -h some-redis