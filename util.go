package fiberpaginate

import "strings"

/* This file contains all the utility and helper functions that can be reused in paginate */

// isAllowedSort returns true if the given sort string is allowed in the given list of allowedSorts,
// false otherwise.
func isAllowedSort(sort string, allowedSorts []string) bool {
	for _, s := range allowedSorts {
		if s == sort {
			return true
		}
	}
	return false
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
	var sortFields []SortField

	for _, field := range fields {
		order := ASC
		if strings.HasPrefix(field, "-") {
			order = DESC
			field = field[1:]
		}

		if isAllowedSort(field, allowedSorts) {
			sortFields = append(sortFields, SortField{Field: field, Order: order})
		}
	}

	if len(sortFields) == 0 {
		return []SortField{{Field: defaultSort, Order: ASC}}
	}

	return sortFields
}
