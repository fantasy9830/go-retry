package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	goretry "github.com/fantasy9830/go-retry"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))

	ctx := context.Background()

	options := []goretry.OptionFunc{
		goretry.MaxRetries(10),
		goretry.WithBackoff(goretry.BackoffLinear(time.Second)),
	}

	goretry.Retry(ctx, func() error {
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			return err
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		log.Println(string(body))

		return err
	}, options...)
}
