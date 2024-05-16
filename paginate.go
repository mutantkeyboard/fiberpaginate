package fiberpaginate

import (
	"github.com/gofiber/fiber/v2"
)

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey byte

// The keys for the values in context
const (
	pageInfoKey contextKey = 0
)

func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		page := c.QueryInt(cfg.PageKey, cfg.DefaultPage)

		limit := c.QueryInt(cfg.LimitKey, cfg.DefaultLimit)

		c.Locals(pageInfoKey, NewPageInfo(page, limit))

		return c.Next()
	}
}

// FromContext returns the PageInfo from the context.
// If there is a PageInfo in the context, it is returned and the boolean is true.
// If there is no PageInfo in the context, nil is returned and the boolean is false.
func FromContext(c *fiber.Ctx) (*PageInfo, bool) {
	if fiberpaginate, ok := c.Locals(pageInfoKey).(*PageInfo); ok {
		return fiberpaginate, true
	}
	return nil, false
}
