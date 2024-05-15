package fiberpaginate

import (
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// DefaultPage is the default page number to use when not provided by the client.
	//
	// Optional. Default: 1
	DefaultPage int

	// DefaultLimit is the default limit to use when not provided by the client.
	//
	// Optional. Default: 10
	DefaultLimit int
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:         nil,
	DefaultPage:  1,
	DefaultLimit: 10,
}

// Helper function to set default values
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
	if cfg.DefaultLimit < 1 {
		cfg.DefaultLimit = ConfigDefault.DefaultLimit
	}
	if cfg.DefaultPage < 1 {
		cfg.DefaultPage = ConfigDefault.DefaultPage
	}

	return cfg
}
