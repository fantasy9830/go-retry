package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	retry "github.com/fantasy9830/go-retry"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))

	options := []retry.OptionFunc{
		retry.WithContext(context.Background()),
		retry.MaxRetries(10),
		retry.WithBackoff(retry.BackoffLinear(time.Second)),
		retry.OnRetry(func(attempt uint, err error) {
			log.Println(attempt, err)
		}),
	}

	retry.Do(func(ctx context.Context) error {
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
