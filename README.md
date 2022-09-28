# go-retry

## Usage
### func MaxRetries
```go
func MaxRetries(maxRetries uint)
```

### func WithBackoff
```go
func WithBackoff(BackoffLinear(duration time.Duration))
func WithBackoff(BackoffExponential(duration time.Duration))
```

## examples
```go
ctx := context.Background()

options := []goretry.OptionFunc{
	goretry.MaxRetries(10),
	goretry.WithBackoff(goretry.BackoffLinear(time.Second)),
}

goretry.Do(ctx, func() error {
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