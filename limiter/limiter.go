package limiter

import "net/http"

// Limiter for the rate of requests to the endpoint
type Limiter struct {
	count int
}

// New limiter
func New() *Limiter {
	return &Limiter{
		count: 0,
	}
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
