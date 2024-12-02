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
go get -u github.com/garrettladley/fiberpaginate/v2
```

## Config

| Property     | Type                    | Description                                                                                                                                                              | Default   |
|--------------|-------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|
| Next         | `func(*fiber.Ctx) bool` | Next defines a function to skip this middleware when returned true.                                                                                                      | `nil`     |
| PageKey      | `string`                | PageKey is the key for the page number in the query string.                                                                                                              | `"page"`  |
| DefaultPage  | `int`                   | DefaultPage is the default page number to use when not provided as a query parameter in the url. If the page number is less than 1, it will be set to the default page number, 1. | `1`       |
| LimitKey     | `string`                | LimitKey is the key for the limit number in the query string.                                                                                                            | `"limit"` |
| DefaultLimit | `int`                   | DefaultLimit is the default limit to use when not provided as a query parameter in the url. If the limit is less than 1, it will be set to the default limit, 10.       | `10`      |
| OffsetKey    | `string`                | OffsetKey is the key for the offset number in the query string.                                                                                                          | `"offset"`|
| SortKey      | `string`                | SortKey is the key for the sort parameter in the query string.                                                                                                           | `"sort"`  |
| DefaultSort  | `string`                | DefaultSort is the default sort to use when not provided as a query parameter in the url.                                                                                | `""`      |
| AllowedSorts | `[]string`              | AllowedSorts is a whitelist of sortable fields. If empty, all fields are allowed.                                                                                        | `nil`     |                                              | `10`                  |


## Example

### Basic Usage with Offset and Sort


```go
package main

import (
	"log"

	"github.com/garrettladley/fiberpaginate/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiberpaginate.New(fiberpaginate.Config{
		PageKey:      "page",
		LimitKey:     "limit",
		OffsetKey:    "offset",
		SortKey:      "sort",
		DefaultSort:  "id",
		AllowedSorts: []string{"id", "name", "date"},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		// This will handle query strings like:
		// ?page=2&limit=10
		// ?offset=20&limit=10
		// ?sort=name,-date
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(fiber.Map{
			"page":   pageInfo.Page,
			"limit":  pageInfo.Limit,
			"offset": pageInfo.Offset,
			"start":  pageInfo.Start(),
			"sort":   pageInfo.Sort,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
```

## With Custom config

Here, you can use our Config to fine tune values for the middlware

```go
package main

import (
	"log"

	"github.com/garrettladley/fiberpaginate/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiberpaginate.New(fiberpaginate.Config{
		PageKey:      "p",
		LimitKey:     "l",
		OffsetKey:    "o",
		SortKey:      "s",
		DefaultPage:  1,
		DefaultLimit: 20,
		DefaultSort:  "created_at",
		AllowedSorts: []string{"id", "name", "created_at"},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		// This will handle query strings like:
		// ?p=2&l=10
		// ?o=40&l=20
		// ?s=name,-created_at
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(fiber.Map{
			"page":   pageInfo.Page,
			"limit":  pageInfo.Limit,
			"offset": pageInfo.Offset,
			"start":  pageInfo.Start(),
			"sort":   pageInfo.Sort,
		})
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

	"github.com/garrettladley/fiberpaginate/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiberpaginate.New(fiberpaginate.Config{
		DefaultPage:  -1,
		DefaultLimit: -1,
		OffsetKey:    "offset",
		SortKey:      "sort",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(fiber.Map{
			"page":   pageInfo.Page,   // Will be 1
			"limit":  pageInfo.Limit,  // Will be 10
			"offset": pageInfo.Offset,
			"start":  pageInfo.Start(),
			"sort":   pageInfo.Sort,
		})
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

	"github.com/garrettladley/fiberpaginate/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiberpaginate.New(fiberpaginate.Config{
		OffsetKey: "offset",
		SortKey:   "sort",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		// When given a query string like ?page=foo&limit=bar&offset=baz&sort=qux,
		// the middleware will parse the query string and set 
		// the values to 0 or default values due to the invalid types
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		return c.JSON(fiber.Map{
			"page":   pageInfo.Page,   // Will be 1 (default)
			"limit":  pageInfo.Limit,  // Will be 10 (default)
			"offset": pageInfo.Offset, // Will be 0
			"start":  pageInfo.Start(),
			"sort":   pageInfo.Sort,   // Will be empty or default
		})
	})

	log.Fatal(app.Listen(":3000"))
}
```
### Generating Next Page URL

This example shows how to use the `NextPageUrl` method to generate a URL for the next page of results as implemented in the `page_info.go` file.

```go
package main

import (
	"log"

	"github.com/garrettladley/fiberpaginate/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(fiberpaginate.New(fiberpaginate.Config{
		PageKey:      "page",
		LimitKey:     "limit",
		SortKey:      "sort",
		DefaultSort:  "id",
		AllowedSorts: []string{"id", "name", "date"},
	}))

	app.Get("/users", func(c *fiber.Ctx) error {
		pageInfo, ok := fiberpaginate.FromContext(c)
		if !ok {
			return fiber.ErrBadRequest
		}

		// Assume we have a function that fetches users based on pageInfo
		// users := fetchUsers(pageInfo)

		nextPageURL := pageInfo.NextPageURL("/users")

		return c.JSON(fiber.Map{
			"page":        pageInfo.Page,
			"limit":       pageInfo.Limit,
			"start":       pageInfo.Start(),
			"sort":        pageInfo.Sort,
			"nextPageURL": nextPageURL,
			// "users":       users,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
