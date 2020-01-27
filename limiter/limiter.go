package limiter

import (
	"fmt"
	"net/http"

	redis "github.com/go-redis/redis/v7"
)

// Limiter for the rate of requests to the endpoint
type Limiter struct {
	count  int
	client *redis.Client
}

// New limiter
func New() (*Limiter, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // YOLO for now
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	fmt.Println(pong, err)

	return &Limiter{
		count: 0,
	}, nil
}

// Limit the rate of requests to this service
func (l *Limiter) Limit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if l.count > 5 {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		l.count++

		f(w, r)
	}
}
