package main

import (
	"io"
	"log"
	"net/http"
)

func rateLimit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if true == true {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		f(w, r)
	}
}

func test1(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from test1!\n")
}

func main() {
	http.HandleFunc("/test1", rateLimit(test1))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
