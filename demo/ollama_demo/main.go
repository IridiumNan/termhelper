package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ollama/ollama/api"
)

func main() {
	// client, err := api.ClientFromEnvironment()
	client := api.NewClient(
		&url.URL{
			Scheme: "http",
			Host:   "100.120.81.67:11434",
		}, &http.Client{
			Timeout: 60 * time.Second,
		},
	)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	req := &api.GenerateRequest{
		Model:  "qwen3.5:4b",
		Prompt: "how many planets are there?",

		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Println(resp.Response)
		return nil
	}

	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}
