package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wwkeyboard/distributed-rate-limiter/limiter"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Hello from %v\n", r.URL.Path)
	io.WriteString(w, msg)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rl, err := limiter.New(3)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/test1", rl.Limit(testHandler))
	http.HandleFunc("/test2", rl.Limit(testHandler))
	http.HandleFunc("/unlimited", testHandler)

	addr := fmt.Sprintf(":%v", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
