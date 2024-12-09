package fiberpaginate

import (
	"github.com/gofiber/fiber/v2"
)

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey struct{}

// The keys for the values in context
var pageInfoKey = contextKey{}

// MaxLimit is the maximum limit allowed which prevents excesive memory usage
const MaxLimit = 100

func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)
	if cfg.DefaultSort == "" {
		cfg.DefaultSort = "id"
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		page := c.QueryInt(cfg.PageKey, cfg.DefaultPage)

		limit := c.QueryInt(cfg.LimitKey, cfg.DefaultLimit)
		if limit > MaxLimit {
			limit = MaxLimit
		}

		offset := c.QueryInt("offset", 0)

		sorts := parseSortQuery(c.Query(cfg.SortKey), cfg.AllowedSorts, cfg.DefaultSort)

		c.Locals(pageInfoKey, NewPageInfo(page, limit, offset, sorts))

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
