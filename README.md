# httpify

`httpify` is a Go package designed to simplify and enhance HTTP client functionality by providing customizable retry strategies and transport configurations. It aims to improve the handling of HTTP requests with flexible retry mechanisms and efficient connection management.


## Installation

To use `httpify` in your Go project, you can install it using `go get`:

```sh
go get github.com/cyinnove/httpify
```

## Usage

Here are some examples of how to use the package:

### Creating a Retry Strategy

```go
package main

import (
	"fmt"
	"time"
	"github.com/cyinnove/httpify"
)

func main() {
	retryStrategy := httpify.DefaultRetryStrategy()
	
	// Example usage of retryStrategy
	fmt.Println(retryStrategy(1*time.Second, 10*time.Second, 2, nil))
}
```

### Creating an HTTP Client with Custom Transport

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/cyinnove/httpify"
)

func main() {
	client := &http.Client{
		Transport: httpify.NoKeepAliveTransport(),
	}
	
	// Example usage of client
	resp, err := client.Get("https://example.com")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Println("Response Status:", resp.Status)
}
```


## Inspiration

This package is inspired by the following projects:
- [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp): A library for retrying HTTP requests with configurable retry strategies.
- [retryablehttp-go](https://github.com/projectdiscovery/retryablehttp-go): Another library for retrying HTTP requests with customizable retry policies.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests to help improve `httpify`.

