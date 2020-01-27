package limiter

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLimiter_Limit(t *testing.T) {
	rl := New()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	})
	ts := httptest.NewServer(rl.Limit(handler))
	defer ts.Close()

	// Test the first request isn't limited
	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == 429 {
		log.Fatal("didn't allow the first request")
	}

	// increment the counter
	http.Get(ts.URL)
	http.Get(ts.URL)
	http.Get(ts.URL)
	http.Get(ts.URL)
	http.Get(ts.URL)

	// test the 5th is limited
	res, err = http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 429 {
		log.Fatal("didn't block the 6th request")
	}
}
