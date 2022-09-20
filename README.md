# go-retry

## examples
```go
	ctx := context.Background()

	options := []goretry.OptionFunc{
		goretry.MaxRetries(10),
		goretry.WithBackoff(goretry.BackoffLinear(time.Second)),
	}

	goretry.Retry(ctx, func() error {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
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

```