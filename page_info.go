package fiberpaginate

// PageInfo contains the pagination information.
type PageInfo struct {
	// Page is the current page number.
	Page int
	// Limit is the number of items per page.
	Limit int
}

func NewPageInfo(page int, limit int) *PageInfo {
	return &PageInfo{
		Page:  page,
		Limit: limit,
	}
}

// Start returns the start index based on the current page and limit.
func (p *PageInfo) Start() int {
	return (p.Page - 1) * p.Limit
}
