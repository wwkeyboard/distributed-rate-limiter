#!/bin/sh

# Starts a development redis server and forwards port 6379

docker run --rm -p 6379:6379 --name some-redis -d redis