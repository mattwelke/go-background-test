package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const (
	GCP_PROJECT_ID               = "GCP_PROJECT_ID"
	GCP_KEY_FILE_BASE64          = "GCP_KEY_FILE_BASE64"
	GCP_PUBSUB_SUBSCRIPTION_NAME = "GCP_PUBSUB_SUBSCRIPTION_NAME"
)

func main() {
	// Print memory stats as it runs
	// go func() {
	// 	timeStart := time.Now().UTC()
	// 	for {
	// 		time.Sleep(time.Duration(60 * time.Second))
	// 		diffStart := time.Now().UTC().Sub(timeStart)
	// 		fmt.Printf("time since startup: %s\n", diffStart)
	// 		PrintMemUsage()
	// 	}
	// }()

	// Also, listen to Pub/Sub
	go func() {
		gcpProjectID := os.Getenv(GCP_PROJECT_ID)
		if gcpProjectID == "" {
			log.Fatalf("missing env var %s", GCP_PROJECT_ID)
		}
		client, err := pubsub.NewClient(context.Background(), gcpProjectID,
			option.WithCredentialsJSON(gcpCredsJSON()))
		if err != nil {
			log.Fatalf("could not create GCP Pub/Sub client: %v", err)
		}
		subName := os.Getenv(GCP_PUBSUB_SUBSCRIPTION_NAME)
		if subName == "" {
			log.Fatalf("missing env var %s", GCP_PUBSUB_SUBSCRIPTION_NAME)
		}
		sub := client.Subscription(subName)
		if err = sub.Receive(context.Background(), func(_ context.Context, m *pubsub.Message) {
			i, _ := strconv.ParseInt(string(m.Data), 10, 64)
			processPubSubMessage(i)
			m.Ack()
		}); err != nil {
			log.Fatalf("could not listen to GCP Pub/Sub subscription: %v", err)
		}
	}()

	// Also, run an HTTP server to convince platform that process is healthy
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// https://golangcode.com/print-the-current-memory-usage/
// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Decodes and returns the JSON creds for the GCP service account from the env var.
func gcpCredsJSON() []byte {
	gcpKeyFileBase64 := os.Getenv(GCP_KEY_FILE_BASE64)
	if gcpKeyFileBase64 == "" {
		log.Fatalf("missing env var %s", GCP_KEY_FILE_BASE64)
	}
	gcpKeyFileDecoded, _ := b64.StdEncoding.DecodeString(gcpKeyFileBase64)
	return gcpKeyFileDecoded
}

// Processes a message number (i) by sleeping for i milliseconds.
func processPubSubMessage(i int64) {
	time.Sleep(time.Duration(i) * time.Millisecond)
	fmt.Printf("Processed message with i = %d.", i)
}
