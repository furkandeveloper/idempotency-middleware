<p align="center">
  <img src="https://github.com/user-attachments/assets/449d7dd6-ae6a-4c52-92ef-bdbe962f1215" style="max-width:100%;" height="300" />
</p>

## Give a Star 🌟
If you liked the project or if **idempotency-middleware** helped you, please give a star.

## Overview

The `idempotency-middleware` library adds idempotency support to your Go API projects. It ensures that a client request is processed only once, regardless of how many times it is received. This middleware is especially useful for handling retry logic in distributed systems.

## Features

- Easy-to-configure idempotency mechanism using Redis.

- Customizable options such as header key and expiration time.

- Seamless integration with the Echo framework.

## Requirements

- Go 1.23.4 or later

- Redis

- Echo v4

## Installation

Add the library to your project using `go get`:

```bash
go get github.com/furkandeveloper/idempotency-middleware
```

Ensure the dependencies in your `go.mod` file include:


```go
module your-module-name

require (
    github.com/furkandeveloper/idempotency-middleware latest
)
```

## Usage
### Configuration

Use the `LoadConfig` function to load the configuration for Redis and the idempotency middleware.

```go
cfg := pkg.LoadConfig(1, "X-Request-Id")
```

**Expiration Time**: Specify the idempotency key expiration time in minutes.

**Header Key**: Define the HTTP header used to identify the idempotency key.

### Initialize Redis Client

Create a Redis client using the loaded configuration:

```go
redisClient := pkg.NewClient(cfg.Redis)
```

### Middleware Setup

Create an instance of the idempotency middleware:

```go
idempotencyMiddleware := pkg.NewIdempotencyMiddleware(redisClient, pkg.Option(cfg.Idempotency))
```

### Apply Midlleware To Routes

```go
package main

import (
	"github.com/furkandeveloper/idempotency-middleware/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ExampleHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "Request processed successfully!"})
}

func ExampleHandlerWithoutIdempotency(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "No idempotency applied."})
}

func main() {
	cfg := pkg.LoadConfig(1, "X-Request-Id")

	redisClient := pkg.NewClient(cfg.Redis)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	idempotencyMiddleware := pkg.NewIdempotencyMiddleware(redisClient, pkg.Option(cfg.Idempotency))

	e.GET("/no-middleware", ExampleHandlerWithoutIdempotency)
	e.GET("/with-middleware", ExampleHandler, idempotencyMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}

```

### How It Works

- Request Header: The client sends a unique token in the header (e.g., X-Request-Id).
- Redis Check: The middleware checks Redis for an existing response for the given token.
- Process:
    -  If a response exists, it is returned directly.
    - If not, the request is processed, and the response is cached in Redis.
- Cache Expiry: Responses are stored in Redis for the configured expiration time.
