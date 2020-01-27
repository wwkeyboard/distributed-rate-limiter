package limiter

import (
	"fmt"
	"net/http"
	"time"

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

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Limiter{
		count:  0,
		client: client,
	}, nil
}

// Limit the rate of requests to this service
func (l *Limiter) Limit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expiration time.Duration
		key := fmt.Sprintf("%v%v", time.Now().Minute(), "/something")

		count, err := l.client.Get(key).Int()
		if err == redis.Nil {
			// the key doesn't exist
			count = 0
			expiration = 1 * time.Minute
		} else if err != nil {
			fmt.Println("here")
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Printf("count for %v -> %v\n", key, count)

		if count > 5 {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		count++
		err = l.client.Set(key, count, expiration).Err()
		if err != nil {
			// log the error but we probably don't need to stop the request from firing
			fmt.Println(err)
		}
		l.count++

		f(w, r)
	}
}
