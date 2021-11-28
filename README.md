# serverfx

serverfx provides a simple configurable HTTP Server implementation with graceful shutdown.

## Getting Started

### Install the package

```shell
go get -u github.com/clarke94/serverfx
```

### Initialise and Serve

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

### Provide Options when initialising to configure the defaults

```go
package foo

import (
	"net/http"
	
	"github.com/clarke94/serverfx"
)

type Handler struct {}

func (h Handler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}

func main() {
    server := serverfx.New(
		Handler{},
		serverfx.WithAddress(":8080"),
		serverfx.WithMaxHeaderBytes(1 << 20),
	)
	
	if err := server.Serve(); err != nil {
		// handle error
    }
}
```
