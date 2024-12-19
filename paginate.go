package fiberpaginate

import (
	"slices"
	"strings"

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
		if page < 1 {
			page = 1
		}

		limit := c.QueryInt(cfg.LimitKey, cfg.DefaultLimit)
		if limit < 1 {
			limit = cfg.DefaultLimit // This will set it to 10 if it's negative or zero
		}
		if limit > MaxLimit {
			limit = MaxLimit
		}

		offset := c.QueryInt("offset", 0)
		if offset < 0 {
			offset = 0
		}

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

// parseSortQuery takes a query string and a list of allowed sorts, and returns a slice of SortFields.
// If the query string is empty, it returns a slice with a single SortField with the given defaultSort.
// The query string is split on commas, and each field is checked against the allowedSorts.
// If a field is not allowed, it is skipped.
// The order of each field is determined by its prefix, with "-" indicating DESC and no prefix indicating ASC.
// If no allowed fields are found, the same single-element slice is returned with the defaultSort.
func parseSortQuery(query string, allowedSorts []string, defaultSort string) []SortField {
	if query == "" {
		return []SortField{{Field: defaultSort, Order: ASC}}
	}

	fields := strings.Split(query, ",")
	sortFields := make([]SortField, 0, len(fields))

	for _, field := range fields {
		order := ASC
		if strings.HasPrefix(field, "-") {
			order = DESC
			field = field[1:]
		}
		if slices.Contains(allowedSorts, field) {
			sortFields = append(sortFields, SortField{Field: field, Order: order})
		}

	}

	if len(sortFields) == 0 {
		return []SortField{{Field: defaultSort, Order: ASC}}
	}

	return sortFields
}
