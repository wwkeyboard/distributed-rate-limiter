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

Development Requirments
=======================

1. docker,
2. golang, tested with 1.13
3. redis-cli, if you want to poke Redis directly
