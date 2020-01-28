# distributed-rate-limiter

This is a demonstration of distributed rate limiting on a go
service. This controls access to a specific endpoint across multiple
nodes. So the path /a would get a given rate allocation, and /b/c
would have its own rate limit.

This coordinated is done through Redis. The first version uses the
rate limiter pattern #1 in the Redis
documentation. https://redis.io/commands/incr#pattern-rate-limiter-1

This pattern uses a single Redis key for each minute & path
combination. So minute 4 of the hour will have 4/first_path and
4/second_path etc.. These keys have a count of requests as their
value.

When a new request comes into the server the middleware first builds
the key based on the path of the request and the minute of the hour,
then pulls the value associated with that key from Redis. The value is
used to check that the number of requests for that path during this
minute isn't above the limit, then the middleware either increments
the counter and allows the request, or blocks the request with a 429
HTTP status code.

Gotchas
=====
- If the service can't reach Redis it sends 500s back to the
  user. This could be changed to allow unrestricted access when Redis
  is unavailable but I feel that in most cases of a degraded situation
  it's better to default to restricting access then allowing it
  unrestricted.
- It requires two calls to Redis, one to check the number of requests
  this minute and one to increment that value. This introduces latency
  into the original request to this service. This latency could be
  reduced by doing the second call to Redis in parallel to servicing
  the request itself. Once we know the request is allowed it's
  acceptable to start servicing it. If updating the counter in Redis
  fails we don't necessarily need to stop access to the limited
  resource. We should, however, log that the update failed and alert
  if the rate of failed updates to Redis becomes too high.


Development
=========

## Requirements

1. docker,
2. golang, tested with 1.13
3. redis-cli, if you want to poke Redis directly

## Running the Server

1. start Redis

``` shell
$ scripts/start-redis-server.sh
```

2. start the server

``` shell
$ PORT=8080 go run main.go
```

3. start another server on another port

``` shell
$ PORT=8081 go run main.go
```

4. test!

``` shell
$ http :8080/test1
$ http :8081/test1
$ http :8080/test1
$ http :8081/test1
$ http :8080/test1
$ http :8081/test1
```

And you should get back a 429

You can also set the request limit of the server by setting the
`LIMIT` env var. This defaults to 100, but it's useful for development
to set it to something much lower. Be careful, if you start two
servers with different limits the one with the lower limit will block
requests without incriminating the total limit count. This will cause
weirdness and in a real production system the limits should come from
a configuration file instead of the environment.

Deployment
========

This requires Redis to be running for the rate limiting to
work. Because this is a critical part of the infrastructure it'll need
to be monitored and run like a production service.

There is no state in Redis older than a minute, so Redis doesn't need
backups. The biggest risk in restarting the service is that all of the
counts are lost, so it would be possible for a spike in load to the
service if Redis is restarted. The rate limiter will also throw 500s
while Redis is down instead of passing all traffic unlimited, it
wouldn't be hard to change this behaviour to fail open instead of
failing closed.

Scaling
=====

Scaling this to large numbers of requests will depend on what the
request pattern looks like. If the requests are spread over a large
number of paths, e.g. 10,000 paths doing 100 requests a minute the
buckets could be sharded over a number of different Redis
instances. For example `/foo/bar` would be handled by instance 1 and
`/foo/baz` by instance 2. The configuration of this would be tedious,
but an endpoint could be migrated from one instance to another in a
very short amount of time since the rate limiter's state is reset
every minute.

Scaling a single path's request load would require more
creativity. The middleware could cache the request count locally, and
only update the central service at a fixed interval, say every 10
requests. If the updates are done out of band with the initial web
request the latency impact to the web request would be minimal. But
there exists the potential to overshoot the limit by (number of
servers * depth of cache) requests.

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
- Use set for the first time a key is set then use the atomic
  increment for future sets. Right now it's done as one uniform action
  for expediency. Right now there is a race condition that means some
  requests could go uncounted. If we use increment instead there is
  still a race condition at the final test, increment, then lock test,
  but it's a much smaller window.
