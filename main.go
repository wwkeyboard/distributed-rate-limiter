package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/wwkeyboard/distributed-rate-limiter/limiter"
)

func test1(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from test1!\n")
}

func main() {
	rl, err := limiter.New(100)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/test1", rl.Limit(test1))
	http.HandleFunc("/unlimited", test1)

	addr := fmt.Sprintf(":%v", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, nil))
}
