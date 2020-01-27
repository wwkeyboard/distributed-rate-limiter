package main

import (
	"io"
	"log"
	"net/http"

	"github.com/wwkeyboard/distributed-rate-limiter/limiter"
)

func test1(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from test1!\n")
}

func main() {
	rl, err := limiter.New()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/test1", rl.Limit(test1))
	http.HandleFunc("/unlimited", test1)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
