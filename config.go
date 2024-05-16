package fiberpaginate

import (
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for the pagination middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// PageKey is the key for the page number in the query string.
	//
	// Optional. Default: "page"
	PageKey string

	// DefaultPage is the default page number to use when not provided by the client.
	// If the page number is less than 1, it will be set to the default page number, 1.
	//
	// Optional. Default: 1
	DefaultPage int

	// LimitKey is the key for the limit number in the query string.
	//
	// Optional. Default: "limit"
	LimitKey string

	// DefaultLimit is the default limit to use when not provided by the client.
	// If the limit is less than 1, it will be set to the default limit, 10.
	//
	// Optional. Default: 10
	DefaultLimit int
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:         nil,
	PageKey:      "page",
	DefaultPage:  1,
	LimitKey:     "limit",
	DefaultLimit: 10,
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.PageKey == "" {
		cfg.PageKey = ConfigDefault.PageKey
	}
	if cfg.DefaultLimit < 1 {
		cfg.DefaultLimit = ConfigDefault.DefaultLimit
	}
	if cfg.LimitKey == "" {
		cfg.LimitKey = ConfigDefault.LimitKey
	}
	if cfg.DefaultPage < 1 {
		cfg.DefaultPage = ConfigDefault.DefaultPage
	}

	return cfg
}
