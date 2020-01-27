// +build integration

package limiter

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
 * This test requires redis running locally,
 * scripts/integration-test.sh sets this up and runs the test.
 */

func TestLimiter_Limit(t *testing.T) {
	rl, err := New(5)
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
