package limiter

import (
	"fmt"
	"net/http"
	"time"

	redis "github.com/go-redis/redis/v7"
)

// Limiter for the rate of requests to the endpoint
type Limiter struct {
	limit  int
	client *redis.Client
}

// New limiter configured for a redis running on localhost
// returns an error if it can't connect to the redis instance
// TODO pass in the redis client instead of building it here
func New(limit int, client *redis.Client) (*Limiter, error) {
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Limiter{
		limit:  limit,
		client: client,
	}, nil
}

// Limit the rate of requests to this service
func (l *Limiter) Limit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expiration time.Duration
		key := key(r.URL.Path)

		count, err := l.client.Get(key).Int()
		if err == redis.Nil {
			// the key doesn't exist
			count = 0
			expiration = 1 * time.Minute
		} else if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// TODO replace this with a logging system, it's printing to stdout for
		// demonstration/debugging
		fmt.Printf("count for %v -> %v\n", key, count)

		if count > l.limit {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		count++
		err = l.client.Set(key, count, expiration).Err()
		if err != nil {
			// log the error but we probably don't need to stop the request from firing
			fmt.Println(err)
		}

		f(w, r)
	}
}

// key is the key of the counting bucket for this endpoint for this minute
func key(slug string) string {
	return fmt.Sprintf("%v%v", time.Now().Minute(), slug)
}
