# distributed-rate-limiter

This is a demonstration of distributed rate limiting on a go
service. This controls access to a specific endpoint across multiple
nodes. So /a would get the given rate allocation, /b/c would have its
own rate limit.

This is coordinated through Redis. The first version uses the rate
limiter pattern #1 in the Redis
documentation. https://redis.io/commands/incr#pattern-rate-limiter-1

TODO:
- [ ] benhhmark to demonstrate limiting (wrk?)

Gotchas
=====
- If the service can't reach Redis it 500s back to the user.
- It requires two calls to Redis, one to check the number of requests
  this minute and one to increment that value. This could be reduced
  to one by doing the call to increment the value in parallel.

Development
=========

## Requirements

1. docker,
2. golang, tested with 1.13
3. redis-cli, if you want to poke Redis directly

## Running the Server

1. start Redis

    scripts/start-redis-server.sh

2. start the server

    PORT=8080 go run main.go

3. start another server on another port

    PORT=8081 go run main.go

4. test!

    http :8080/test1
    http :8081/test1
    http :8080/test1
    http :8081/test1
    http :8080/test1
    http :8081/test1

And you should get back a 429

You can set the request limit of the server by setting the LIMIT env
var. This defaults to 100, but it's useful for development to set it
to something much lower. Be careful, if you start two servers with
different limits the one with the lower limit will block requests
without incriminating the total limit count. This will cause weirdness
and in a real production system the limits should come from a
configuration file instead of the environment.

Deployment
========

This requires Redis to be running for the rate limiting to
work. Because this is a critical part of the infrastructure it'll need
to be monitored and run like a production service.

There is no state in Redis older than a minute, so it doesn't need
backups. The biggest risk in restarting the service is that all of the
counts are lost, so it would be possible for a spike in load to the
service if Redis is restarted. The rate limiter will also throw 500s
while Redis is down instead of passing all traffic unlimited, it
wouldn't be hard to change this behaviour to fail open instead of
failing closed.

Future Work
========

- Configuration, right now all endpoints have the same cap of 100
  requests.
- Metrics, a process could scrape all of the keys in Redis and present
  them to the metrics system(or operate as a Prometheus scrape). This
  would allow the construction of a graph of what endpoints have the
  most traffic. It would also allow a rough estimate of the load on
  each endpoint. This estimate could be made more accurate by
  incriminating the counter past the limit.
- Logging, right now it's printing errors and status for each request
  to stdout. This should log in the same manner as the other systems
  it's deployed with.
