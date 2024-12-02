package fiberpaginate

import "fmt"

// SortOrder represents sort order
type SortOrder string

// SortOrder constants
const (
	ASC  SortOrder = "asc"
	DESC SortOrder = "desc"
)

// SortField represents sort field
type SortField struct {
	Field string
	Order SortOrder
}

// SortOrderFromString returns a SortOrder from a string
func SortOrderFromString(s string) SortOrder {
	switch s {
	case "asc":
		return ASC
	case "desc":
		return DESC
	default:
		return ASC
	}
}

// PageInfo contains the pagination information.
type PageInfo struct {
	// Page is the current page number.
	Page int
	// Limit is the number of items per page.
	Limit int
	// Offset is the offset of the current page.
	Offset int
	// Sort is the sort order.
	Sort []SortField
}

func NewPageInfo(page int, limit int, offset int, sort []SortField) *PageInfo {
	return &PageInfo{
		Page:   page,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
}

// Start returns the start index based on the current page and limit.
func (p *PageInfo) Start() int {
	if p.Offset > 0 {
		return p.Offset
	}
	return (p.Page - 1) * p.Limit
}

// SortBy adds a new sort field to the sort order.
func (p *PageInfo) SortBy(field string, order SortOrder) *PageInfo {
	p.Sort = append(p.Sort, SortField{Field: field, Order: order})
	return p
}

// NextPageURL returns the URL for the next page given the baseURL.
//
// For example, if baseURL is "https://example.com/users" and the current page is 1
// with a limit of 10, the returned URL would be
// "https://example.com/users?page=2&limit=10".
func (p *PageInfo) NextPageURL(baseURL string) string {
	return fmt.Sprintf("%s?page=%d&limit=%d", baseURL, p.Page+1, p.Limit)
}
