# distributed-rate-limiter

This is a demonstration of distributed rate limiting on a go
service. This controls access to a specific endpoint across multiple
nodes. So /a would get the given rate allocation, /b/c would have its
own rate limit.

This is coordinated through Redis. The first version uses the rate
limiter pattern #1 in the Redis
documentation. https://redis.io/commands/incr#pattern-rate-limiter-1

Things to consider:
- [ ] configuration
- [ ] metrics/status
- [ ] deployment/operation of Redis + the service

TODO:
- [ ] basic web service
- [ ] go server
- [ ] middleware
- [ ] benhhmark to demonstrate limiting (wrk?)

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

And you should see a 429 be returned!
