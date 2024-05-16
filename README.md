<h1 align="center">fiberpaginate</h1>

<div align="center">
  <a href="https://pkg.go.dev/github.com/garrettladley/fiberpaginate#section-documentation">
    <img src="https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white"
      alt="Go.Dev Reference" />
  </a>
  <a href="https://goreportcard.com/report/github.com/garrettladley/fiberpaginate">
    <img src="https://goreportcard.com/badge/github.com/garrettladley/fiberpaginate"
      alt="fiberpaginate Go Report" />
  </a>
  <a href="https://github.com/garrettladley/fiberpaginate/actions/workflows/test.yml">
    <img src="https://github.com/garrettladley/fiberpaginate/actions/workflows/test.yml/badge.svg"
      alt="Test Workflow Status" />
  </a>  
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/license-MIT-brightgreen.svg"
      alt="MIT License" />
  </a>
</div>

<br/>

<div align="center">
  <strong>A pagination middleware for <a href="https://gofiber.io">Fiber</a></strong>
</div>

## Install

```console
go get -u github.com/garrettladley/fiberpaginate
```

## Config

| Property            | Type                        | Description                                                                                                                   | Default                |
|:--------------------|:----------------------------|:------------------------------------------------------------------------------------------------------------------------------|:-----------------------|
| Next              | `func(*fiber.Ctx) bool`     | Next defines a function to skip this middleware when returned true.                                                                                     | `nil`                  |
| PageKey              | `string`     | PageKey is the key for the page number in the query. string                                                                                     | `"page"`                  |
| DefaultPage    | `int`             | DefaultPage is the default page number to use when not provided by the client.                                                   | `1`       |
| LimitKey              | `string`     | LimitKey is the key for the limit number in the query. string                                                                                     | `"limit"`                  |
| DefaultLimit        | `int`                  | DefaultLimit is the default limit to use when not provided by the client.                                                                   | `10`                  |

## Example

```go
package main

import (
	"log"

	"github.com/garrettladley/fiberpaginate"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiber_pagnate.New())

	app.Get("/", func(c *fiber.Ctx) error {
		// when given a query string like ?page=2&limit=5,
		// the middleware will parse the query string and set the values
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(
			fiber.Map{
				"page":   pageInfo.Page,
				"limit":  pageInfo.Limit,
				"start": pageInfo.Start(),
			},
		)
	})

	log.Fatal(app.Listen(":3000"))
}
```

## Note with negative values in the config

If DefaultPage is configured to be less than 1, the middleware will use 1 as the default value for the page.
If DefaultLimit is configured to be less than 1, the middleware will use 10 as the default value for the limit.

```go
package main

import (
	"log"

	"github.com/garrettladley/fiberpaginate"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiberpaginate.New(fiberpaginate.Config{
		DefaultPage:  -1,
		DefaultLimit: -1,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(
			fiber.Map{
				"page":   pageInfo.Page, // 1
				"limit":  pageInfo.Limit, // 10
				"start": pageInfo.Start(), // 0
			},
		)
	})

	log.Fatal(app.Listen(":3000"))
}
```

## Note with invalid keys

If the client provides an invalid type to page or limit, the middleware will use 0 as the value stored for the page or limit.

```go
package main

import (
	"log"

	"github.com/garrettladley/fiberpaginate"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiberpaginate.New())

	app.Get("/", func(c *fiber.Ctx) error {
		// when given a query string like ?page=foo&limit=bar,
		// the middleware will parse the query string and set 
		// the values to 0 due to the invalid types
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(
			fiber.Map{
				"page":   pageInfo.Page, // 0
				"limit":  pageInfo.Limit, // 0
				"start": pageInfo.Start(), // 0
			},
		)
	})

	log.Fatal(app.Listen(":3000"))
}
```
