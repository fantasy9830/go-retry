# go-retry

## Installation
```shell
go get github.com/fantasy9830/go-retry
```

## Usage
### func WithContext
```go
func WithContext(ctx context.Context)
```

### func MaxRetries
```go
func MaxRetries(maxRetries uint)
```

### func WithBackoff
```go
func WithBackoff(BackoffLinear(duration time.Duration))
func WithBackoff(BackoffExponential(duration time.Duration))
```

### func OnRetry
```go
func OnRetry(func(attempt uint, err error) {})
```

## examples
```go
options := []retry.OptionFunc{
	retry.WithContext(context.Background()),
	retry.MaxRetries(10),
	retry.WithBackoff(retry.BackoffLinear(time.Second)),
	retry.OnRetry(func(attempt uint, err error) {
		log.Println(attempt, err)
	}),
}

retry.Do(func(ctx context.Context) error {
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