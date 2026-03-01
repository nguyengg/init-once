# Golang sync.Once improvements for intialisers

[![Go Reference](https://pkg.go.dev/badge/github.com/nguyengg/init-once.svg)](https://pkg.go.dev/github.com/nguyengg/init-once)

Get with:
```shell
go get github.com/nguyengg/init-once
```

`init.Once` makes sure the initialiser is called exactly once regardless of its return status:
```go
package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	init "github.com/nguyengg/init-once"
)

type App struct {
	Client *s3.Client
	once   init.Once
}

func (a *App) DoSomething(ctx context.Context) error {
	// because init.Once is used, subsequent a.once.Do will always return the same error, nil or non-nil.
	if err := a.once.Do(func() error {
		return a.init(ctx)
	}); err != nil {
		return err
	}

	// there's no risk of a nil Client pointer here due to the err check above.
	// sync.Once with an extra initErr or sync.OnceValue would have been needed to achieve the same thing.
	_, _ = a.Client.GetObject(ctx, &s3.GetObjectInput{})

	return nil
}

func (a *App) init(ctx context.Context) error {
	if a.Client == nil {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}

		a.Client = s3.NewFromConfig(cfg)
	}

	return nil
}

```

If you want to be able to retry initilisation over and over until first success, use `init.SuccessOnce`:
```go
package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	init "github.com/nguyengg/init-once"
)

type App struct {
	Client *s3.Client
	once   init.SuccessOnce
}

func (a *App) DoSomething(ctx context.Context) error {
	// because init.SuccessOnce is used, init can be retried until its first success!
	if err := a.once.Do(func() error {
		return a.init(ctx)
	}); err != nil {
		return err
	}

	// there's still no risk of a nil Client pointer here due to the err check above.
	_, _ = a.Client.GetObject(ctx, &s3.GetObjectInput{})

	return nil
}

func (a *App) init(ctx context.Context) error {
	if a.Client == nil {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}

		a.Client = s3.NewFromConfig(cfg)
	}

	return nil
}

```
