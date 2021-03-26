package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	go func() {
		timeStart := time.Now().UTC()
		for {
			time.Sleep(time.Duration(1 * time.Second))
			diffStart := time.Now().UTC().Sub(timeStart)
			fmt.Printf("time since startup: %s\n", diffStart)
		}
	}()

	// Server just to keep platform happy
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
