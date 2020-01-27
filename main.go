package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v7"
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

	var err error
	limit := 100
	envLimit := os.Getenv("LIMIT")
	if envLimit != "" {
		limit, err = strconv.Atoi(envLimit)
		if err != nil {
			log.Fatalf("Failed to parse requested limit %v", err)
		}
	}

	rl, err := limiter.New(
		limit,
		redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // YOLO for now
			DB:       0,
		}))
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/test1", rl.Limit(testHandler))
	http.HandleFunc("/test2", rl.Limit(testHandler))
	http.HandleFunc("/unlimited", testHandler)

	addr := fmt.Sprintf(":%v", port)
	fmt.Println("Listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
