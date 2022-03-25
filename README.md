# serverfx

[![CircleCI](https://circleci.com/gh/clarke94/serverfx/tree/main.svg?style=svg)](https://circleci.com/gh/clarke94/serverfx/tree/main)
[![Go Reference](https://pkg.go.dev/badge/github.com/clarke94/serverfx.svg)](https://pkg.go.dev/github.com/clarke94/serverfx)


serverfx provides a simple configurable HTTP Server implementation with graceful shutdown.

## Getting Started

### Install the package

```shell
go get -u github.com/clarke94/serverfx
```

### Usage

### Initialising and Serving a HTTP server.

```go
package foo

import (
	"net/http"
	
	"github.com/clarke94/serverfx"
)

type Handler struct {}

func (h Handler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func main() {
    server := serverfx.New(Handler{})
	
	if err := server.Serve(); err != nil { 
	    // handle error 
	}
}
```

### Providing Options to configure the defaults

```go
package foo

import (
	"net/http"
	"time"

	"github.com/clarke94/serverfx"
)

type Handler struct{}

func (h Handler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func main() {
	server := serverfx.New(
		Handler{},
		serverfx.WithAddress(":8080"),
		serverfx.WithMaxHeaderBytes(1<<20),
		serverfx.WithGracefulTimeout(10*time.Second),
	)

	if err := server.Serve(); err != nil {
	    // handle error
	}
}
```

### Manually calling Shutdown

```go
package foo

import (
	"errors"
	"net/http"
	"time"

	"github.com/clarke94/serverfx"
)

type Handler struct{}

func (h Handler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func main() {
	server := serverfx.New(
		Handler{},
		serverfx.WithAddress(":8080"),
		serverfx.WithMaxHeaderBytes(1<<20),
		serverfx.WithGracefulTimeout(10*time.Second),
	)

	// after 5 seconds manually shutdown the server.
	go func() {
		time.Sleep(5 * time.Second)
		server.Shutdown()
	}()

	if err := server.Serve(); err != nil {
		// The errors from Shutdown can be handled 
		// here once the Server has Shutdown.
		if errors.Is(err, serverfx.ErrUnableToGracefulShutdown) {}

		if errors.Is(err, serverfx.ErrUnableToListenAndServe) {}
	}
}
```

### Generic example with the Gin framework

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/clarke94/serverfx"
)

func main() {
	router := gin.New()
	server := serverfx.New[*gin.Engine](router)

	server.Handler.GET("/foo", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	if err := server.Serve(); err != nil {
		// handle error
	}
}
```

## License

This project uses the [MIT License](LICENSE).
